package models

import (
	"time"

	"github.com/gocql/gocql"
)

type EventRaw struct {
	ProjectID    gocql.UUID        `cql:"project_id"`
	EventName    string            `cql:"event_name"`
	EventTime    time.Time         `cql:"event_time"`
	InsertedTime time.Time         `cql:"inserted_time"`
	Payload      map[string]string `cql:"payload"`
}
