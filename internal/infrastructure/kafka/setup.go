package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/SUT-technology/log-analysis/internal/infrastructure/cassandra"
	clickhouse "github.com/SUT-technology/log-analysis/internal/infrastructure/clickHouse"
	"github.com/segmentio/kafka-go"
)

// KafkaClient wraps producer and consumer.
type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

// NewKafkaProducer initializes a Kafka writer.
func NewKafkaProducer(brokers []string, topic string) *KafkaClient {
	return &KafkaClient{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// NewKafkaConsumer initializes a Kafka reader.
func NewKafkaConsumer(brokers []string, topic, groupID string) *KafkaClient {
	return &KafkaClient{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

// Produce sends a message to Kafka.
func (k *KafkaClient) Produce(ctx context.Context, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   nil,
		Value: b,
		Time:  time.Now(),
	})
}

// Consume reads a message from Kafka.
func (k *KafkaClient) Consume(ctx context.Context) ([]byte, error) {
	m, err := k.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return m.Value, nil
}

// ProcessAndInsert consumes messages from Kafka and inserts them into Cassandra and ClickHouse.
func (k *KafkaClient) ProcessAndInsert(ctx context.Context, cass *cassandra.CassandraClient, ch *clickhouse.ClickHouseSQLClient) error {
	for {
		msg, err := k.Consume(ctx)
		if err != nil {
			log.Printf("Error consuming message: %v", err)
			continue
		}

		// Deserialize the message into a LogMessage
		var logMessage models.LogMessage
		if err := json.Unmarshal(msg, &logMessage); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Insert into Cassandra
		eventRaw := models.EventRaw{
			ProjectID:    logMessage.ProjectID,
			EventName:    logMessage.Name,
			EventTime:    logMessage.Timestamp,
			InsertedTime: logMessage.Timestamp,
			Payload:      logMessage.Payload,
		}
		if err := cass.InsertEvent(ctx, &eventRaw, 3600); err != nil {
			log.Printf("Error inserting into Cassandra: %v", err)
		}

		// Insert into ClickHouse
		if err := ch.InsertEvent(ctx, logMessage); err != nil {
			log.Printf("Error inserting into ClickHouse: %v", err)
		}
	}
}
