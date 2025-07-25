package rest

import (
	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/interface/handler/authhndlr"
	"github.com/SUT-technology/log-analysis/internal/interface/handler/loghndlr"
	"github.com/SUT-technology/log-analysis/internal/interface/handler/projecthndlr"
	"github.com/labstack/echo/v4"
)

func register(e *echo.Echo, srvc application.Services, m *middlewares) {
	// Create groups with middleware
	// swaggerGroup := NewGroup("/swagger", middlewares.loggingMiddleware, mux

	logs := e.Group("/api/logs", m.corsMiddleware(), m.JWTMiddleware)
	loghndlr.New(logs, srvc)
	projects := e.Group("/api/projects", m.corsMiddleware(), m.JWTMiddleware)
	projecthndlr.New(projects, srvc)

	auth := e.Group("/api/auth")
	authhndlr.New(auth, srvc)

}
