package models

import "time"

type LogMessage struct {
	ProjectID string            `json:"project_id"`
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   map[string]string `json:"payload"`
}
