package dto

import "time"

type EventSummary struct {
	EventName  string    `json:"event_name"`
	LastOccur  time.Time `json:"last_occur"`
	TotalCount int64     `json:"total_count"`
}

// ListEventsResponse پاسخ لیست خلاصه ایونت‌ها
type ListEventsResponse struct {
	ProjectID string         `json:"project_id"`
	Filters   EventFilters   `json:"applied_filters"`
	Data      []EventSummary `json:"data"`
}

// EventDetail شامل جزئیات یک ایونت
type EventDetail struct {
	EventName    string            `json:"event_name"`
	EventTime    time.Time         `json:"event_time"`
	InsertedTime time.Time         `json:"inserted_time"`
	Payload      map[string]string `json:"payload"`
}

// DetailEventsResponse پاسخ جزئیات و ناوبری
type DetailEventsResponse struct {
	ProjectID   string       `json:"project_id"`
	Filters     EventFilters `json:"applied_filters"`
	Current     EventDetail  `json:"current_event"`
	PrevEventID string       `json:"prev_event_id,omitempty"`
	NextEventID string       `json:"next_event_id,omitempty"`
}

// EventFilters ورودی فیلتر
type EventFilters struct {
	ProjectID      string            `json:"project_id" query:"project_id" binding:"required,uuid"`
	SearchableKeys map[string]string `json:"searchable_keys,omitempty"`
	EventName      string            `json:"event_name,omitempty" query:"event_name"`
	Page           int               `json:"page,omitempty" query:"page"`
}

