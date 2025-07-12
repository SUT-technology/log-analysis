package application

import (
	"github.com/SUT-technology/log-analysis/internal/application/logsrvc"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/cassandra"
	clickhouse "github.com/SUT-technology/log-analysis/internal/infrastructure/clickHouse"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/kafka"
)

type Services struct {
	LogSrvc logsrvc.LogSrvc
}

func New(producer *kafka.KafkaClient,
	consumer *kafka.KafkaClient,
	cassandra *cassandra.CassandraClient,
	cockroachdb *cockroachdb.CockroachDBClient,
	clickhouse *clickhouse.ClickHouseSQLClient) Services {

	return Services{
		LogSrvc: logsrvc.New(producer, cassandra, cockroachdb, clickhouse),
	}
}
