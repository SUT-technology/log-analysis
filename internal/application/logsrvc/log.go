package logsrvc

import (
	"context"
	"fmt"

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
		fmt.Errorf("error debug: %w", err)
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

func (s LogSrvc) ListEvents(ctx context.Context, filters dto.EventFilters) (dto.ListEventsResponse, error) {
	const pageSize = 10
	offset := (filters.Page - 1) * pageSize
	log.Info("offset: ", offset)

	fmt.Printf("[debug] projectID: %s", filters.ProjectID)

	whereClause := fmt.Sprintf("WHERE project_id = '%s'", filters.ProjectID)

	if filters.EventName != "" {
		whereClause += fmt.Sprintf(" AND event_name = '%s'", filters.EventName)
	}

	var i = 1
	for key, value := range filters.SearchableKeys {
		if value != "" {
			whereClause += fmt.Sprintf(" AND payload.key[%v] = '%s' AND payload.value[%v] = '%s'", i, key, i, value)
		}
		i++
	}

	query := fmt.Sprintf(`SELECT event_name, max(event_time) AS last_occur, count(*) AS total_count FROM events_clickhouse %v GROUP BY event_name ORDER BY last_occur DESC LIMIT %d OFFSET %d`, whereClause, pageSize, offset)

	fmt.Printf("[debug] query: %s", query)

	rows, err := s.clickhouse.DB.QueryContext(ctx, query)
	if err != nil {
		log.Error("error in getting rows:", err)
		return dto.ListEventsResponse{}, err
	}
	defer rows.Close()

	var events []dto.EventSummary
	for rows.Next() {
		var e dto.EventSummary
		if err := rows.Scan(&e.EventName, &e.LastOccur, &e.TotalCount); err != nil {
			log.Error("error in getting scanning:", err)
			return dto.ListEventsResponse{}, err
		}
		events = append(events, e)
	}

	return dto.ListEventsResponse{
		ProjectID: filters.ProjectID,
		Data:      events,
		Filters:   filters,
	}, nil
}

func (s LogSrvc) DetailEvent(ctx context.Context, filters dto.EventFilters) (dto.DetailEventsResponse, error) {

	event, err := s.cassandra.GetEventByTime(ctx, filters.ProjectID, filters.EventTime, filters.EventName)

	if err != nil {
		log.Error("error getting event by time from cassandra: ", err)
		return dto.DetailEventsResponse{}, err
	}

	nextTime := s.clickhouse.GetNextEventTime(ctx, filters)

	prevTime := s.clickhouse.GetPreviousEventTime(ctx, filters)

	return dto.DetailEventsResponse{
		ProjectID: filters.ProjectID,
		Filters:   filters,
		Current: dto.EventDetail{
			EventName:    event.EventName,
			EventTime:    event.EventTime,
			InsertedTime: event.InsertedTime,
			Payload:      event.Payload,
		},
		NextEventTime:     nextTime,
		PreviousEventTime: prevTime,
	}, nil
}
