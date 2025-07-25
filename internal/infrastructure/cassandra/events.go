package cassandra

import (
	"context"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

func (c *CassandraClient) InsertEvent(ctx context.Context, evt *models.EventRaw, ttl int) error {
	return c.session.Query(
		`INSERT INTO events_raw (project_id, event_name, event_time, inserted_time, payload) VALUES (?, ?, ?, ?, ?) USING TTL ?`,
		evt.ProjectID.String(), evt.EventName, evt.EventTime, time.Now(), evt.Payload, ttl,
	).WithContext(ctx).Exec()
}

func (c *CassandraClient) GetEventByTime(ctx context.Context, projectID string, eventTime time.Time, eventName string) (*models.EventRaw, error) {
	query := `SELECT project_id, event_name, event_time, inserted_time, payload 
			  FROM events_raw 
			  WHERE project_id = ? AND event_time = ?`

	var event models.EventRaw
	var id string

	log.Info("[debug] query: ", query)
	err := c.session.Query(query, projectID, eventTime).WithContext(ctx).Scan(
		&id,
		&event.EventName,
		&event.EventTime,
		&event.InsertedTime,
		&event.Payload,
	)

	if err != nil {
		return nil, err
	}
	event.ProjectID = uuid.MustParse(id)
	return &event, nil
}
