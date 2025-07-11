package cockroachdb

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
)

func (c *CockroachDBClient) GetUser(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	row := c.db.QueryRowContext(ctx, `SELECT id, username, password_hash, created_at FROM users WHERE id = $1`, id)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}
