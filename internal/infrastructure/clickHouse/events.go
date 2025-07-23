package clickhouse

import (
	"context"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
)

// InsertEvent writes raw event into ClickHouse with TTL.
func (c *ClickHouseSQLClient) InsertEvent(ctx context.Context, logMessage models.LogMessage) error {
	// Convert payload map to string array
	payload := make([]string, 0, len(logMessage.Payload)*2)
	for k, v := range logMessage.Payload {
		payload = append(payload, k, v)
	}

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events_clickhouse 
		(project_id, event_name, event_time, inserted_time, payload) 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		logMessage.ProjectID,
		logMessage.Name,
		logMessage.Timestamp,
		time.Now(), // inserted_time
		payload,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

