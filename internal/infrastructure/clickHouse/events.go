package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/SUT-technology/log-analysis/internal/domain/models"
)

func (c *ClickHouseSQLClient) InsertEvent(ctx context.Context, logMessage models.LogMessage) error {
	keys := make([]string, 0, len(logMessage.Payload))
	values := make([]string, 0, len(logMessage.Payload))
	for k, v := range logMessage.Payload {
		keys = append(keys, k)
		values = append(values, v)
	}

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin error: %w", err)
	}
	defer tx.Rollback()

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
		keys,
		values,
	)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func (c *ClickHouseSQLClient) GetNextEventTime(ctx context.Context, filters dto.EventFilters) time.Time {
	conditions := []string{fmt.Sprintf("project_id = '%s'", filters.ProjectID)}

	if filters.EventName != "" {
		conditions = append(conditions, fmt.Sprintf("event_name = '%s'", filters.EventName))
	}

	var i = 1
	for key, value := range filters.SearchableKeys {
		conditions = append(conditions, fmt.Sprintf("payload.key[%v] = '%s'", i, key))
		conditions = append(conditions, fmt.Sprintf("payload.value[%v] = '%s'", i, value))
		i++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT event_time FROM events_clickhouse %s AND event_time > '%s' ORDER BY event_time ASC LIMIT 1`, whereClause, filters.EventTime.Format("2006-01-02 15:04:05"))

	row := c.DB.QueryRowContext(ctx, query)
	var next time.Time
	if err := row.Scan(&next); err != nil {
		return time.Time{}
	}
	return next
}

func (c *ClickHouseSQLClient) GetPreviousEventTime(ctx context.Context, filters dto.EventFilters) time.Time {
	conditions := []string{fmt.Sprintf("project_id = '%s'", filters.ProjectID)}

	if filters.EventName != "" {
		conditions = append(conditions, fmt.Sprintf("event_name = '%s'", filters.EventName))
	}

	var i = 1
	for key, value := range filters.SearchableKeys {
		conditions = append(conditions, fmt.Sprintf("payload.key[%v] = '%s'", i, key))
		conditions = append(conditions, fmt.Sprintf("payload.value[%v] = '%s'", i, value))
		i++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT event_time FROM events_clickhouse %s AND event_time < '%s' ORDER BY event_time DESC LIMIT 1`, whereClause, filters.EventTime.Format("2006-01-02 15:04:05"))

	row := c.DB.QueryRowContext(ctx, query)
	var prev time.Time
	if err := row.Scan(&prev); err != nil {
		return time.Time{}
	}
	return prev
}
