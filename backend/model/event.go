package model

import "time"

// KeyValue represents a dynamic key-value pair in a log event
type KeyValue struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type Event struct {
    ID        string     `json:"id"`         // UUID
    ProjectID string     `json:"project_id"` // related project
    Name      string     `json:"name"`       // event name
    Timestamp time.Time  `json:"timestamp"`  // when the event occurred
    Payload   []KeyValue `json:"payload"`    // full event details
    CreatedAt time.Time  `json:"created_at"` // when it was ingested
}
