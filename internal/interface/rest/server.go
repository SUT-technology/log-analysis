package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/interface/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Server struct {
	srv    *echo.Echo
	defers []func()
}

func NewServer(srvc application.Services, cfg config.Config) *Server {
	fmt.Println("🟢 NewServer called")
	var dfrs []func()

	e := echo.New()
	e.Debug = true
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &Validator{validator: validator.New()}

	closer := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := e.Shutdown(ctx)
		if err != nil {
			slog.Error("close HTTP server", err)
		}
	}
	dfrs = append(dfrs, closer)

	// manage middlewares
	var middleware []echo.MiddlewareFunc
	m := newMiddlewares(cfg)

	middleware = append(middleware, m.corsMiddleware())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Error("Unhandled error:", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	e.Use(middleware...)

	register(e, srvc, m)

	return &Server{srv: e, defers: dfrs}
}

func (s *Server) Start(addr string) error {
	return s.srv.Start(addr)
}

func (s *Server) Stop() {
	for _, f := range slices.Backward(s.defers) {
		f()
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, fieldErr := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("%s failed on '%s'", fieldErr.Field(), fieldErr.Tag()))
			}
			return errors.New(strings.Join(errMsgs, ", "))
		}
		return errors.New("validation failed with unknown error")
	}
	return nil
}
