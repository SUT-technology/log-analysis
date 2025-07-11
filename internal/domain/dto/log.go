package dto

import "time"

type SendLogRequest struct {
	ProjectID string            `json:"project_id" binding:"required,uuid"`
	APIKey    string            `json:"api_key" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	Timestamp time.Time         `json:"timestamp" binding:"required"`
	Payload   map[string]string `json:"payload" binding:"required"`
}

type SendLogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
