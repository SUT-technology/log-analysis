package cockroachdb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type CockroachDBClient struct {
	db *sql.DB
}

func NewCockroachDBClient(dsn string) (*CockroachDBClient, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	c := &CockroachDBClient{db: db}
	c.InitCockroachSchema()
	return c, nil
}

func (c *CockroachDBClient) InitCockroachSchema() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username STRING NOT NULL UNIQUE,
			password_hash STRING NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
			name STRING NOT NULL,
			api_key STRING NOT NULL UNIQUE,
			searchable_keys ARRAY<STRING> NOT NULL,
			ttl_seconds INT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);`,
	}
	for _, stmt := range stmts {
		if _, err := c.db.Exec(stmt); err != nil {
			return fmt.Errorf("cockroach init error: %w", err)
		}
	}
	return nil
}
