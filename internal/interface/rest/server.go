package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/SUT-technology/log-analysis/internal/application"
	"github.com/SUT-technology/log-analysis/internal/interface/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Server struct {
	srv    *echo.Echo
	defers []func()
}

func NewServer(srvc application.Services, cfg config.Config) *Server {
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

	// if cfg.Server.Logger {
	// 	middleware = append(middleware, m.loggerMiddleware)
	// }

	// application specific middlewares
	middleware = append(middleware, m.corsMiddleware())

	// default recover middleware
	// middleware = append(middleware, m.recoverMiddleware)

	// applying middlewares and create a new server
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
