package command

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/cassandra"
	clickhouse "github.com/SUT-technology/log-analysis/internal/infrastructure/clickHouse"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/kafka"
	"github.com/SUT-technology/log-analysis/internal/interface/config"
	"github.com/SUT-technology/log-analysis/internal/interface/rest"
)

func Run() error {
	var configPath string
	flag.StringVar(&configPath, "cfg", "assets/config.yaml", "Configuration File")
	flag.Parse()
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	// راه‌اندازی CockroachDB
	crdb, err := cockroachdb.NewCockroachDBClient(cfg.Cockroach.DSN)
	if err != nil {
		log.Fatalf("cockroach init: %v", err)
	}
	fmt.Println("Connected to CockroachDB")

	// راه‌اندازی Kafka
	producer := kafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	consumer := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	fmt.Println("Kafka producer and consumer ready")

	// راه‌اندازی Cassandra
	cass, err := cassandra.NewCassandraClient(cfg.Cassandra.Hosts, cfg.Cassandra.Keyspace)
	if err != nil {
		log.Fatalf("cassandra init: %v", err)
	}
	fmt.Println("Connected to Cassandra")

	// راه‌اندازی ClickHouse
	chDsn := fmt.Sprintf(
		"tcp://%s?username=%s&password=%s&database=%s",
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Username,
		cfg.ClickHouse.Password,
		cfg.ClickHouse.Database,
	)
	fmt.Println("ClickHouse DSN:", chDsn)
	clickhouseClient, err := clickhouse.NewClickHouseSQLClient(chDsn)
	if err != nil {
		log.Fatalf("clickhouse init: %v", err)
	}
	fmt.Println("Connected to ClickHouse2")

	srvc := application.New(producer, consumer, cass, crdb, clickhouseClient)

	// می‌توانید از این کلاینت‌ها در سرویستان استفاده کنید...
	_ = crdb
	_ = producer
	_ = consumer
	_ = cass
	_ = clickhouseClient
	_ = context.Background()
	fmt.Printf("cfg.Server.Addr = %#v\n", cfg.Server.Addr)
	httpSrv := rest.NewServer(srvc, cfg)
	defer func() {
		slog.Debug("gracefully stopping HTTP server")
		httpSrv.Stop()
	}()

	startErr := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("❌ Server panicked:", r)
			}
		}()
		fmt.Println("🚀 Starting HTTP server on", cfg.Server.Addr)
		err := httpSrv.Start(cfg.Server.Addr)
		startErr <- fmt.Errorf("HTTP server startup: %w", err)
	}()

	go func() {
		log.Println("Starting Kafka consumer...")
		if err := consumer.ProcessAndInsert(context.Background(), cass, clickhouseClient); err != nil {
			log.Fatalf("Error processing Kafka messages: %v", err)
		}
	}()

	select {
	case err := <-startErr:
		slog.Error("failed to start server", slog.Any("error", err))
		return err
	case <-quit:
		slog.Info("received signal to stop server")
		return nil
	}

}
