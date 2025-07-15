package projectsrvc

import (
	"context"

	"github.com/SUT-technology/log-analysis/internal/domain/dto"
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

func (p ProjectSrvc) ProjectsList(ctx context.Context, userID string) (dto.ProjectsListRespone,error) {
	projects,err:= p.cockroachdb.GetProjects(ctx,userID)
	if err != nil {
		return dto.ProjectsListRespone{},err
	}
	var projectSummeries = make([]dto.ProjectSummery,10)
	for _,project := range projects {
		summery := dto.ProjectSummery {
			ProjectID: project.ID,
			ProjectName: project.Name,
			SearchableKeys: project.SearchableKeys,
		}
		projectSummeries = append(projectSummeries, summery)
	}

	return dto.ProjectsListRespone{
		UserID: uuid.MustParse(userID),
		Projects: projectSummeries,
	},nil
}