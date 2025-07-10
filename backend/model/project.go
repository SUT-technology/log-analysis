package model

import "time"

type Project struct {
    ID         string        `json:"id"`          // UUID
	UserID	   string		 `json:"userId"`	  // UUID
    Name       string        `json:"name"`
    SearchKeys []string      `json:"search_keys"` // searchable keys
    ApiKey     string        `json:"api_key"`     // secure project key
    TTL        time.Duration `json:"ttl"`         // how long to retain logs
    CreatedAt  time.Time     `json:"created_at"`
}