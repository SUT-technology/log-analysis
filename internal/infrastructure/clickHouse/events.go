package clickhouse

import (
	"context"
	"github.com/SUT-technology/log-analysis/internal/domain/models"
)


// InsertEvent writes raw event into ClickHouse with TTL.
func (c *ClickHouseSQLClient) InsertEvent(ctx context.Context, logMessage models.LogMessage) error {
	query := `
		INSERT INTO events_clickhouse (project_id, event_name, event_time, inserted_time, payload)
		VALUES (?, ?, ?, now(), ?)
	`
	payload := make([]string, 0, len(logMessage.Payload))
	for k, v := range logMessage.Payload {
		payload = append(payload, k, v)
	}

	_, err := c.DB.ExecContext(ctx, query,
		logMessage.ProjectID,
		logMessage.Name,
		logMessage.Timestamp,
		payload,
	)
	return err
}
