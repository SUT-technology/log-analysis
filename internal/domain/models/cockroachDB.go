package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}

type Project struct {
	ID             uuid.UUID `db:"id"`
	OwnerID        uuid.UUID `db:"owner_id"`
	Name           string    `db:"name"`
	APIKey         string    `db:"api_key"`
	SearchableKeys []string  `db:"searchable_keys"`
	TTL            int       `db:"ttl_seconds"`
	CreatedAt      time.Time `db:"created_at"`
}
