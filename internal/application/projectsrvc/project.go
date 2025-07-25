package projectsrvc

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/SUT-technology/log-analysis/internal/domain/models"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/google/uuid"
)

type ProjectSrvc struct {
	cockroachdb *cockroachdb.CockroachDBClient
}

func New(cockroachdb *cockroachdb.CockroachDBClient) ProjectSrvc {
	return ProjectSrvc{
		cockroachdb: cockroachdb,
	}
}

func (p ProjectSrvc) ProjectsList(ctx context.Context, userID uuid.UUID) (dto.ProjectsListRespone, error) {
	projects, err := p.cockroachdb.GetProjects(ctx, userID)
	if err != nil {
		return dto.ProjectsListRespone{}, err
	}
	var projectSummeries []dto.ProjectSummery
	for _, project := range projects {
		summery := dto.ProjectSummery{
			ProjectID:      project.ID,
			ProjectName:    project.Name,
			SearchableKeys: project.SearchableKeys,
		}
		projectSummeries = append(projectSummeries, summery)
	}

	return dto.ProjectsListRespone{
		UserID:   userID,
		Projects: projectSummeries,
	}, nil
}

func (p ProjectSrvc) SaveProject(ctx context.Context, project dto.NewProjectRequset) (dto.NewProjectResponse, error) {

	var model = models.Project{
		OwnerID:        project.OwnerID,
		Name:           project.Name,
		APIKey:         project.APIKey,
		SearchableKeys: project.SearchableKeys,
		TTL:            project.TTL,
	}

	if err := p.cockroachdb.InsertProject(ctx, &model); err != nil {
		return dto.NewProjectResponse{}, err
	}

	return dto.NewProjectResponse{
		ID:             (model.ID).String(),
		OwnerID:        (model.OwnerID).String(),
		Name:           model.Name,
		APIKey:         model.APIKey,
		SearchableKeys: model.SearchableKeys,
		TTL:            model.TTL,
	}, nil
}
