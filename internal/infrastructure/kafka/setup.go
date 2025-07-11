package kafka

import (
	"context"
	"encoding/json"
	"time"

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
