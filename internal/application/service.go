package application

import (
	"github.com/SUT-technology/log-analysis/internal/application/authsrvc"
	"github.com/SUT-technology/log-analysis/internal/application/logsrvc"
	"github.com/SUT-technology/log-analysis/internal/application/projectsrvc"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/cassandra"
	clickhouse "github.com/SUT-technology/log-analysis/internal/infrastructure/clickHouse"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/kafka"
	"github.com/SUT-technology/log-analysis/internal/interface/config"
)

type Services struct {
	LogSrvc logsrvc.LogSrvc
	AuthSrvc    authsrvc.AuthSrvc
	ProjectSrvc projectsrvc.ProjectSrvc
}

func New(producer *kafka.KafkaClient,
	consumer *kafka.KafkaClient,
	cassandra *cassandra.CassandraClient,
	cockroachdb *cockroachdb.CockroachDBClient,
	clickhouse *clickhouse.ClickHouseSQLClient,
	cfg config.Config) Services {

	return Services{
		ProjectSrvc: projectsrvc.New(cockroachdb),
		LogSrvc:     logsrvc.New(producer, cassandra, cockroachdb, clickhouse),
		AuthSrvc:    authsrvc.NewAuthSrvc(cockroachdb, cfg.Server.SecretKey),
	}
}
