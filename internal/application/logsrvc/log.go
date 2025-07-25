package logsrvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/cassandra"
	clickhouse "github.com/SUT-technology/log-analysis/internal/infrastructure/clickHouse"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/kafka"
	"github.com/labstack/gommon/log"
)

type LogSrvc struct {
	producer    *kafka.KafkaClient
	cassandra   *cassandra.CassandraClient
	cockroachdb *cockroachdb.CockroachDBClient
	clickhouse  *clickhouse.ClickHouseSQLClient
}

func New(producer *kafka.KafkaClient,
	cassandra *cassandra.CassandraClient,
	cockroachdb *cockroachdb.CockroachDBClient,
	clickhouse *clickhouse.ClickHouseSQLClient) LogSrvc {
	return LogSrvc{
		producer:    producer,
		cassandra:   cassandra,
		cockroachdb: cockroachdb,
		clickhouse:  clickhouse,
	}
}

func (s LogSrvc) SendLog(ctx context.Context, req dto.SendLogRequest) (dto.SendLogResponse, error) {
	project, err := s.cockroachdb.GetProject(ctx, req.ProjectID)
	if err != nil {
		fmt.Errorf("error debug: %w",err)
		return dto.SendLogResponse{}, err
	}

	if project == nil {
		return dto.SendLogResponse{}, fmt.Errorf("project not found")
	}

	if project.APIKey != req.APIKey {
		return dto.SendLogResponse{}, fmt.Errorf("invalid API key")
	}

	kafkaDto := models.LogMessage{
		ProjectID: project.ID,
		Name:      req.Name,
		Timestamp: req.Timestamp,
		Payload:   req.Payload,
	}

	err = s.producer.Produce(ctx, kafkaDto)
	if err != nil {
		return dto.SendLogResponse{}, err
	}

	return dto.SendLogResponse{}, nil
}

// ListEvents لیستی از خلاصه ایونت‌ها با فیلتر
func (s LogSrvc) ListEvents(ctx context.Context, filters dto.EventFilters) (dto.ListEventsResponse, error) {
	const pageSize = 10
	offset := (filters.Page-1) * pageSize

	fmt.Printf("[debug] projectID: %s",filters.ProjectID)

	whereClause := fmt.Sprintf("WHERE project_id = '%s'",filters.ProjectID)

	if filters.EventName != "" {
		whereClause += fmt.Sprintf(" AND event_name = '%s'",filters.EventName)
	}

	for key,value := range filters.SearchableKeys {
		whereClause += fmt.Sprintf(" AND payload['%s'] = '%s'", key, value)
	}

	query := fmt.Sprintf(`SELECT event_name, max(event_time) AS last_occur, count(*) AS total_count FROM events_clickhouse %v GROUP BY event_name ORDER BY last_occur DESC LIMIT %d OFFSET %d`, whereClause, pageSize, offset)

	fmt.Printf("[debug] query: %s",query)

	// Run query
	rows, err := s.clickhouse.DB.QueryContext(ctx, query)
	if err != nil {
		log.Error("error in getting rows:",err)
		return dto.ListEventsResponse{}, err
	}
	defer rows.Close()

	var events []dto.EventSummary
	for rows.Next() {
		var e dto.EventSummary
		if err := rows.Scan(&e.EventName, &e.LastOccur, &e.TotalCount); err != nil {
			log.Error("error in getting scanning:",err)
			return dto.ListEventsResponse{}, err
		}
		events = append(events, e)
	}

	return dto.ListEventsResponse{
		ProjectID: filters.ProjectID,
		Data: events,
		Filters: filters,
	}, nil
}


// DetailEvent جزئیات ایونت فعلی و ناوبری را بازمی‌گرداند
func (s LogSrvc) DetailEvent(ctx context.Context, filters dto.EventFilters) (dto.DetailEventsResponse, error) {

	// Build WHERE conditions dynamically
	conditions := []string{fmt.Sprintf("project_id = '%s'",filters.ProjectID)}

	if filters.EventName != "" {
		conditions = append(conditions, fmt.Sprintf("event_name = '%s'",filters.EventName))
	}

	for key := range filters.SearchableKeys {
		conditions = append(conditions, fmt.Sprintf("payload['%s'] = '%s'", key, key))
		// params[key] = value
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var query string

	if filters.Position == dto.Absolute {
		query = fmt.Sprintf(`SELECT * FROM events_clickhouse %s AND event_time = '%s'`, whereClause, filters.EventTime.Format("2006-01-02 15:04:05"))
	} else if filters.Position == dto.Next {
		query = fmt.Sprintf(`SELECT * FROM events_clickhouse %s AND event_time > '%s' ORDER BY event_time ASC LIMIT 1`, whereClause, filters.EventTime.Format("2006-01-02 15:04:05"))
	} else if filters.Position == dto.Previous {
		query = fmt.Sprintf(`SELECT * FROM events_clickhouse %s AND event_time < '%s' ORDER BY event_time DESC LIMIT 1`, whereClause, filters.EventTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("[debug] query: %s",query)

	// Run query
	row, err := s.clickhouse.DB.QueryContext(ctx, query)
	if err != nil {
		return dto.DetailEventsResponse{}, err
	}
	defer row.Close()

	var event models.EventRaw
	var payloadKeys []string
	var PayloadValues []string
	if row.Next() {
		if err := row.Scan(&event.ProjectID,&event.EventName, &event.EventTime, &event.InsertedTime, &payloadKeys,&PayloadValues); err != nil {
			return dto.DetailEventsResponse{}, err
		}
	}
	var payload = make(map[string]string)
	for i := 0; i < len(payloadKeys); i++ {
		payload[payloadKeys[i]] = PayloadValues[i]
	}

	return dto.DetailEventsResponse{
		ProjectID: filters.ProjectID,
		Filters: filters,
		Current: dto.EventDetail{
			EventName: event.EventName,
			EventTime: event.EventTime,
			InsertedTime: event.InsertedTime,
			Payload: payload,
		},
	}, nil
}
