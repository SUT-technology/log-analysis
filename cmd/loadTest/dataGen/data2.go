package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var clickhouseDSN = "tcp://127.0.0.1:9000?username=default&password=&database=default"

func main() {
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
	projectIDs = []uuid.UUID{
		uuid.MustParse("a63d3d4f-656f-49f6-be2c-7cab588092f7"),
		uuid.MustParse("2000cc4f-64b3-4760-806f-fdc2910f338a"),
		uuid.MustParse("342a666f-2873-4355-963c-88dd72081c4a"),
		uuid.MustParse("a7ad2ae3-9b29-4e05-b3d0-ca224d746051"),
	}

	eventNames := []string{"login", "logout", "purchase", "view"}
	now := time.Now()

	for _, pid := range projectIDs {
		for i := 0; i < 5000; i++ {
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
