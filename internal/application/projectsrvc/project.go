package projectsrvc

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/google/uuid"
)

type ProjectSrvc struct {
	cockroachdb *cockroachdb.CockroachDBClient
}

func NewProjectSrvc(cockroachdb *cockroachdb.CockroachDBClient) ProjectSrvc {
	return ProjectSrvc{
		cockroachdb: cockroachdb,
	}
}

func (c ProjectSrvc) GetUserProjects(ctx context.Context, userId uuid.UUID) ([]*models.Project, error) {
	projects, err := c.cockroachdb.GetUserProjects(ctx, userId)
	return projects, err
}
