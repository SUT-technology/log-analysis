package logsrvc

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/dto"
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

	return dto.SendLogResponse{}, nil
}

// ListEvents لیستی از خلاصه ایونت‌ها با فیلتر
func (s LogSrvc) ListEvents(ctx context.Context, filters dto.EventFilters) (dto.ListEventsResponse, error) {

	return dto.ListEventsResponse{}, nil
}

// DetailEvent جزئیات ایونت فعلی و ناوبری را بازمی‌گرداند
func (s LogSrvc) DetailEvent(ctx context.Context, filters dto.EventFilters) (dto.DetailEventsResponse, error) {

	return dto.DetailEventsResponse{}, nil
}
