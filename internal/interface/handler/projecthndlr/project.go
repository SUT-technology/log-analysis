package projecthndlr

import (
	"net/http"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ProjectHndlr struct {
	Services application.Services
}

func New(g *echo.Group, srvc application.Services) *ProjectHndlr {
	handler := &ProjectHndlr{Services: srvc}

	g.GET("", handler.GetProjects)

	return handler
}

func (h *ProjectHndlr) GetProjects(c echo.Context) error {
	userId, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil
	}
	resp, err := h.Services.ProjectSrvc.GetUserProjects(c.Request().Context(), userId)
	if err != nil {
		log.Error("Failed to send log:", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp)
}
