package clickhouse

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

type ClickHouseSQLClient struct {
	db *sql.DB
}

// NewClickHouseSQLClient با DSN به ClickHouse وصل می‌شود.
// مثال DSN: "tcp://clickhouse:9000?username=default&password=&database=default"
func NewClickHouseSQLClient(dsn string) (*ClickHouseSQLClient, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("open clickhouse: %w", err)
	}
	// می‌توانی تنظیمات pool را هم اینجا اعمال کنی
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	c := &ClickHouseSQLClient{db: db}
	if err := c.InitClickHouseSchema(); err != nil {
		return nil, fmt.Errorf("init clickhouse schema: %w", err)
	}
	return c, nil
}

func (c *ClickHouseSQLClient) InitClickHouseSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS events_clickhouse (
			project_id UUID,
			event_name String,
			event_time DateTime,
			inserted_time DateTime,
			payload Nested(
				key String,
				value String
			)
		)
		ENGINE = MergeTree()
		PARTITION BY toYYYYMM(event_time)
		ORDER BY (project_id, event_name, event_time);
	`
	if _, err := c.db.Exec(query); err != nil {
		return fmt.Errorf("clickhouse init error: %w", err)
	}
	return nil
}
