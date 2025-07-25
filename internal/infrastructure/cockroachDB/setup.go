package cockroachdb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type CockroachDBClient struct {
	db *sql.DB
}

func NewCockroachDBClient(dsn string) (*CockroachDBClient, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	c := &CockroachDBClient{db: db}
	err = c.InitCockroachSchema()
	if err != nil {
		return nil, fmt.Errorf("cockroachDB init error: %w", err)
	}
	return c, nil
}

func (c *CockroachDBClient) InitCockroachSchema() error {

	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username STRING NOT NULL UNIQUE,
			password_hash STRING NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("cockroach init (users): %w", err)
	}

	_, err = c.db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
			name STRING NOT NULL,
			api_key STRING NOT NULL UNIQUE,
			searchable_keys STRING[] NOT NULL,
			ttl_seconds INT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("cockroach init (projects): %w", err)
	}

	// Insert a test project and user
	// _, err = c.InsertTestUser()
	// if err != nil {
	// 	return fmt.Errorf("insert test user: %w", err)
	// }
	// if err := c.InsertTestProject(ownerID); err != nil {
	// 	return fmt.Errorf("insert test project: %w", err)
	// }
	return nil
}

func (c *CockroachDBClient) InsertTestProject(ownerID string) error {
	stmt := `
	INSERT INTO projects (owner_id, name, api_key, searchable_keys, ttl_seconds)
	VALUES ($1, $2, $3, $4, $5);
	`

	name := "Test Project"
	apiKey := "test-api-key-1274s1a"
	searchableKeys := []string{"key1", "key2", "key3"}
	ttlSeconds := 3600

	_, err := c.db.Exec(stmt, ownerID, name, apiKey, pq.Array(searchableKeys), ttlSeconds)
	if err != nil {
		return fmt.Errorf("failed to insert test project: %w", err)
	}
	return nil
}

func (c *CockroachDBClient) InsertTestUser() (string, error) {
	stmt := `
	INSERT INTO users (username, password_hash)
	VALUES ($1, $2)
	RETURNING id;
	`

	username := "mahdi001"
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	var userID string
	err = c.db.QueryRow(stmt, username, hashedPassword).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to insert test user: %w", err)
	}
	return userID, nil
}
