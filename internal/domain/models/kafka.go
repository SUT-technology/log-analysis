package models

import (
	"time"

	"github.com/google/uuid"
)


type LogMessage struct {
	ProjectID uuid.UUID        `json:"project_id"`
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   map[string]string `json:"payload"`
}
