package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
)

// InsertEvent writes raw event into ClickHouse with TTL.
func (c *ClickHouseSQLClient) InsertEvent(ctx context.Context, logMessage models.LogMessage) error {
	// Convert map to two parallel string slices
	keys := make([]string, 0, len(logMessage.Payload))
	values := make([]string, 0, len(logMessage.Payload))
	for k, v := range logMessage.Payload {
		keys = append(keys, k)
		values = append(values, v)
	}

	// Begin transaction
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin error: %w", err)
	}
	defer tx.Rollback() // safe in case of failure

	// Prepare statement with nested fields
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events_clickhouse 
		(project_id, event_name, event_time, inserted_time, payload.key, payload.value) 
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare error: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		logMessage.ProjectID,
		logMessage.Name,
		logMessage.Timestamp,
		time.Now(),
		keys,   // payload.key
		values, // payload.value
	)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}


