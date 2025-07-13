package models

import (
	"time"

	"github.com/gocql/gocql"
)


type LogMessage struct {
	ProjectID gocql.UUID        `json:"project_id"`
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   map[string]string `json:"payload"`
}
