package db

// import (
//     kafka "github.com/segmentio/kafka-go"
//     "log"
// )

// func EnsureLogTopic() {
//     conn, err := kafka.Dial("tcp", "localhost:9092")
//     if err != nil {
//         log.Fatalf("Failed to connect to Kafka: %v", err)
//     }
//     defer conn.Close()

//     topic := "logs"

//     err = conn.CreateTopics(kafka.TopicConfig{
//         Topic:             topic,
//         NumPartitions:     3,
//         ReplicationFactor: 3,
//     })
//     if err != nil {
//         log.Printf("Topic %s might already exist: %v", topic, err)
//     } else {
//         log.Printf("Topic %s created or exists", topic)
//     }
// }
