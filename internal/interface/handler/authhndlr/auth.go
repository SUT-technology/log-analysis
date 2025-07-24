package authhndlr

import (
	"net/http"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthHndlr struct {
	Services application.Services
}

func New(g *echo.Group, srvc application.Services) *AuthHndlr {
	handler := &AuthHndlr{Services: srvc}

	g.POST("/login", handler.Login)

	return handler
}

func (h *AuthHndlr) Login(c echo.Context) error {
	var req dto.LoginRequest
	log.Info("Received request to login")
	if err := c.Bind(&req); err != nil {
		log.Error("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, dto.LoginResponse{})
	}
	resp, err := h.Services.AuthSrvc.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		log.Error("Failed to send log:", err)
		return c.JSON(http.StatusInternalServerError, dto.LoginResponse{})
	}
	return c.JSON(http.StatusOK, dto.LoginResponse{Token: resp})
}
