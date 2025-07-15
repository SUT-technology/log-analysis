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

	g.GET(":userID",handler.ProjectsList)

	return handler
}

func (p *ProjectHndlr) ProjectsList(c echo.Context) error {
	var filters dto.EventFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ListEventsResponse{Filters: dto.EventFilters{}})
	}
	userID := c.Param("userID")
	resp, err := p.Services.ProjectSrvc.ProjectsList(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ListEventsResponse{Filters: filters})
	}
	return c.JSON(http.StatusOK, resp)
}