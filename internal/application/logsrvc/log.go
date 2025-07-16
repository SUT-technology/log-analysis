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
	offset := filters.Page * pageSize

	// Build WHERE conditions dynamically
	conditions := []string{"project_id = {project_id:String}"}

	if filters.EventName != "" {
		conditions = append(conditions, "event_name = {event_name:String}")
		// params["event_name"] = filters.EventName
	}

	for key := range filters.SearchableKeys {
		conditions = append(conditions, fmt.Sprintf("payload['%s'] = {%s:String}", key, key))
		// params[key] = value
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			event_name,
			max(timestamp) AS last_occur,
			count(*) AS total_count
		FROM events
		%s
		GROUP BY event_name
		ORDER BY last_occur DESC
		LIMIT %d OFFSET %d
	`, whereClause, pageSize, offset)

	// Run query
	rows, err := s.clickhouse.DB.QueryContext(ctx, query)
	if err != nil {
		return dto.ListEventsResponse{}, err
	}
	defer rows.Close()

	var events []dto.EventSummary
	for rows.Next() {
		var e dto.EventSummary
		if err := rows.Scan(&e.EventName, &e.LastOccur, &e.TotalCount); err != nil {
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
	conditions := []string{fmt.Sprintf("project_id = {%s:String}",filters.ProjectID)}

	if filters.EventName != "" {
		conditions = append(conditions, "event_name = {event_name:String}")
	}

	for key := range filters.SearchableKeys {
		conditions = append(conditions, fmt.Sprintf("payload['%s'] = {%s:String}", key, key))
		// params[key] = value
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var query string

	if filters.Position == dto.Absolute {
		query = fmt.Sprintf(`
				SELECT * FROM events
				WHERE %s AND event_time = %s
				GROUP BY event_name
				ORDER BY event_time DESC
				LIMIT 1
				`, whereClause, filters.EventTime)
	} else if filters.Position == dto.Next {
		query = fmt.Sprintf(`
				SELECT * FROM events
				WHERE %s AND event_time > %s
				GROUP BY event_name
				ORDER BY event_time ASC
				LIMIT 1
				`, whereClause, filters.EventTime)
	} else if filters.Position == dto.Previous {
		query = fmt.Sprintf(`
				SELECT * FROM events
				WHERE %s AND event_time < %s
				GROUP BY event_name
				ORDER BY event_time DESC
				LIMIT 1
				`, whereClause, filters.EventTime)
	}

	// Run query
	row, err := s.clickhouse.DB.QueryContext(ctx, query)
	if err != nil {
		return dto.DetailEventsResponse{}, err
	}
	defer row.Close()

	var event dto.EventDetail
	if row.Next() {
		if err := row.Scan(&event.EventName, &event.EventTime, &event.InsertedTime, &event.Payload); err != nil {
			return dto.DetailEventsResponse{}, err
		}
	}

	return dto.DetailEventsResponse{
		ProjectID: filters.ProjectID,
		Filters: filters,
		Current: event,
	}, nil
}
