package cassandra

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
)

// InsertEvent writes raw event into Cassandra with TTL.
func (c *CassandraClient) InsertEvent(ctx context.Context, evt *models.EventRaw, ttl int) error {
	return c.session.Query(
		`INSERT INTO events_raw (project_id, event_name, event_time, inserted_time, payload) VALUES (?, ?, ?, ?, ?) USING TTL ?`,
		evt.ProjectID, evt.EventName, evt.EventTime, evt.InsertedTime, evt.Payload, ttl,
	).WithContext(ctx).Exec()
}
