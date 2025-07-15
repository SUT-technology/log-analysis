package cockroachdb

import (
	"context"
	"strings"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/google/uuid"
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

func (c *CockroachDBClient) GetProjects(ctx context.Context, userID string) ([]models.Project, error) {
	projects := make([]models.Project,10)
	rows,err := c.db.QueryContext(ctx,`SELECT * FROM projects WHERE owner_id = $1`,uuid.MustParse(userID))
	if err != nil {
		return nil,err
	}

	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.OwnerID, &project.Name, &project.APIKey, &project.SearchableKeys, &project.TTL, &project.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects,nil
}
