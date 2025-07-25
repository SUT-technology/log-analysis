package loghndlr

import (
	"fmt"
	"net/http"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/domain/dto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type LogHndlr struct {
	Services application.Services
}

func New(g *echo.Group, srvc application.Services) *LogHndlr {
	handler := &LogHndlr{Services: srvc}

	g.POST("/summery", handler.ListEvents)
	g.POST("/detail", handler.DetailEvent)
	g.POST("", handler.SendLog)

	return handler
}

func (h *LogHndlr) SendLog(c echo.Context) error {
	var req dto.SendLogRequest
	log.Info("Received request to send log")
	if err := c.Bind(&req); err != nil {
		log.Error("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, dto.SendLogResponse{Success: false, Message: "invalid request"})
	}
	resp, err := h.Services.LogSrvc.SendLog(c.Request().Context(), req)
	if err != nil {
		log.Error("Failed to send log:", err)
		return c.JSON(http.StatusInternalServerError, dto.SendLogResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *LogHndlr) ListEvents(c echo.Context) error {
	var filters dto.EventFilters

	log.Info("Received request to summery")

	if err := c.Bind(&filters); err != nil {
		log.Error("Received request to summery err: ", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Info("filters: ", filters)

	// Manually extract searchable_keys[...] from query params
	// filters.SearchableKeys = make(map[string]string)
	// for key, values := range c.QueryParams() {
	// 	if strings.HasPrefix(key, "searchable_keys[") && strings.HasSuffix(key, "]") {
	// 		innerKey := key[len("searchable_keys[") : len(key)-1]
	// 		if len(values) > 0 {
	// 			filters.SearchableKeys[innerKey] = values[0]
	// 		}
	// 	}
	// }

	eventsResponse, err := h.Services.LogSrvc.ListEvents(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, eventsResponse)
}

func (h *LogHndlr) DetailEvent(c echo.Context) error {
	var filters dto.EventFilters
	fmt.Println("Received request for event details:")
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp, err := h.Services.LogSrvc.DetailEvent(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}
