package loghndlr

import (
	"fmt"
	"net/http"
	"strings"

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
	g.GET("/:eventName", handler.DetailEvent)
	g.GET("/:eventName", handler.DetailEvent)
	g.POST("", handler.SendLog)

	return handler
}

// SendLog ورودی را به DTO تبدیل و به سرویس ارسال می‌کند
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

// ListEvents لیست خلاصه‌ایونت‌ها را با فیلتر برمی‌گرداند
func (h *LogHndlr) ListEvents(c echo.Context) error {
	var filters dto.EventFilters

	// Bind simple fields like ProjectID, EventName, etc.
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Manually extract searchable_keys[...] from query params
	filters.SearchableKeys = make(map[string]string)
	for key, values := range c.QueryParams() {
		if strings.HasPrefix(key, "searchable_keys[") && strings.HasSuffix(key, "]") {
			innerKey := key[len("searchable_keys[") : len(key)-1]
			if len(values) > 0 {
				filters.SearchableKeys[innerKey] = values[0]
			}
		}
	}

	eventsResponse, err := h.Services.LogSrvc.ListEvents(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, eventsResponse)
}

// DetailEvent جزئیات یک ایونت خاص را برمی‌گرداند
func (h *LogHndlr) DetailEvent(c echo.Context) error {
	var filters dto.EventFilters
	fmt.Println("Received request for event details:", c.Param("eventName"))
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// eventName از مسیر می‌آید
	filters.EventName = c.Param("eventName")
	resp, err := h.Services.LogSrvc.DetailEvent(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}
