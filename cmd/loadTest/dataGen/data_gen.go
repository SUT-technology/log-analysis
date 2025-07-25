package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	cockroachDSN     = "postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable"
	clickhouseDSN    = "tcp://127.0.0.1:9000?username=default&password=&database=default"
	numUsers         = flag.Int("users", 10, "Number of users to generate")
	projectsPerUser  = flag.Int("projectsPerUser", 5, "Number of projects per user")
	eventsPerProject = flag.Int("eventsPerProject", 1000, "Number of events per project")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// اتصال به CockroachDB
	crDB, err := sql.Open("postgres", cockroachDSN)
	if err != nil {
		log.Fatalf("cockroach connect error: %v", err)
	}
	defer crDB.Close()

	// اتصال به ClickHouse
	chDB, err := sql.Open("clickhouse", clickhouseDSN)
	if err != nil {
		log.Fatalf("clickhouse connect error: %v", err)
	}
	defer chDB.Close()

	// آماده‌سازی statement برای ClickHouse
	tx, err := chDB.Begin()
	if err != nil {
		log.Fatalf("clickhouse tx begin: %v", err)
	}
	stmt, err := tx.Prepare(`
        INSERT INTO events_clickhouse
        (project_id, event_name, event_time, inserted_time, payload.key, payload.value)
        VALUES (?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		log.Fatalf("prepare clickhouse stmt: %v", err)
	}
	defer stmt.Close()

	var projectIDs []uuid.UUID

	// تولید کاربران و پروژه‌ها
	for u := 0; u < *numUsers; u++ {
		username := fmt.Sprintf("user_%d_%s", u, uuid.New().String()[:8])
		passwordHash := uuid.New().String()
		var userID uuid.UUID

		err := crDB.QueryRow(
			`INSERT INTO users(username, password_hash) VALUES($1, $2) RETURNING id`,
			username, passwordHash,
		).Scan(&userID)
		if err != nil {
			log.Fatalf("insert user error: %v", err)
		}

		for p := 0; p < *projectsPerUser; p++ {
			name := fmt.Sprintf("proj_%d_%d", u, p)
			apiKey := uuid.New().String()
			keys := []string{"k1", "k2", "k3"}
			ttl := 3600

			var projectID uuid.UUID
			err := crDB.QueryRow(`
                INSERT INTO projects(owner_id, name, api_key, searchable_keys, ttl_seconds)
                VALUES($1, $2, $3, $4, $5) RETURNING id
            `, userID, name, apiKey, pq.Array(keys), ttl).Scan(&projectID)
			if err != nil {
				log.Fatalf("insert project error: %v", err)
			}

			projectIDs = append(projectIDs, projectID)
		}
	}

	log.Printf("Generated %d projects, starting event insertion...", len(projectIDs))

	// تولید لاگ‌ها
	eventNames := []string{"login", "logout", "purchase", "view"}
	now := time.Now()

	for _, pid := range projectIDs {
		for i := 0; i < *eventsPerProject; i++ {
			evt := eventNames[rand.Intn(len(eventNames))]
			evtTime := now.Add(-time.Duration(rand.Intn(3600*24*30)) * time.Second)
			keys := []string{"k1", "k2", "k3"}
			values := []string{
				fmt.Sprintf("v1_%d", rand.Intn(100)),
				fmt.Sprintf("v2_%d", rand.Intn(100)),
				fmt.Sprintf("v3_%d", rand.Intn(100)),
			}

			if _, err := stmt.Exec(
				pid, evt, evtTime, now,
				keys, values,
			); err != nil {
				log.Fatalf("clickhouse insert error: %v", err)
			}

			if i%10000 == 0 && i > 0 {
				log.Printf("Inserted %d events for project %s", i, pid)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("clickhouse tx commit: %v", err)
	}

	log.Println("Data generation completed.")
}
