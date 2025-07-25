package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProjectSummery struct {
	ProjectID      uuid.UUID `json:"project_id"`
	ProjectName    string    `json:"project_name"`
	SearchableKeys []string  `json:"searchable_keys"`
}

type ProjectsListRespone struct {
	UserID   uuid.UUID        `json:"user_id"`
	Projects []ProjectSummery `json:"projects"`
}

type NewProjectRequset struct {
	OwnerID        uuid.UUID `json:"owner_id"`
	Name           string    `json:"name"`
	APIKey         string    `json:"api_key"`
	SearchableKeys []string  `json:"searchable_keys"`
	TTL            int       `json:"ttl_seconds"`
}

type NewProjectResponse struct {
	ID             string    `json:"id"`
	OwnerID        string    `json:"owner_id"`
	Name           string    `json:"name"`
	APIKey         string    `json:"api_key"`
	SearchableKeys []string  `json:"searchable_keys"`
	TTL            int       `json:"ttl_seconds"`
	CreatedAt      time.Time `json:"created_at"`
}
