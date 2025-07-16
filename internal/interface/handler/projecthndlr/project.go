package projecthndlr

import (
	"net/http"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/labstack/echo/v4"
)

type ProjectHndlr struct {
	Services application.Services
}

func New(g *echo.Group, srvc application.Services) *ProjectHndlr {
	handler := &ProjectHndlr{Services: srvc}


	g.GET("/:userID",handler.ProjectsList)
	g.POST("",handler.CreateProject)

	return handler
}

func (p *ProjectHndlr) ProjectsList(c echo.Context) error {
	userID := c.Param("userID")
	resp, err := p.Services.ProjectSrvc.ProjectsList(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (p *ProjectHndlr) CreateProject(c echo.Context) error {
	var project dto.NewProjectRequset
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp,err := p.Services.ProjectSrvc.SaveProject(c.Request().Context(),project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

