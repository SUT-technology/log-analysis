package cockroachdb

import (
	"context"
	"strings"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/lib/pq"
)

func (c *CockroachDBClient) GetProject(ctx context.Context, id string) (*models.Project, error) {
	id = strings.TrimSpace(id) // ✅ this fixes the error

	project := &models.Project{}
	row := c.db.QueryRowContext(ctx, `
		SELECT id, owner_id, name, api_key, searchable_keys, ttl_seconds, created_at
		FROM projects WHERE id = $1`, id)

	if err := row.Scan(&project.ID, &project.OwnerID, &project.Name, &project.APIKey,
		pq.Array(&project.SearchableKeys), &project.TTL, &project.CreatedAt); err != nil {
		return nil, err
	}
	return project, nil
}
