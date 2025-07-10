package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func InitClickHouse() *sql.DB {
    conn, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return nil
	}

    _,err = conn.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
            project_id UUID,
            event_name String,
            timestamp DateTime,
            key_values Map(String, String)
        ) ENGINE = MergeTree()
        PARTITION BY toYYYYMM(timestamp)
        ORDER BY (project_id, event_name, timestamp)
    `)
	if err != nil {
        log.Fatalf("Failed to create ClickHouse table: %v", err)
    }

    return conn
}
