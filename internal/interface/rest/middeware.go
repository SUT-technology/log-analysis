package rest

import (
	"net/http"
	"strings"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	"github.com/SUT-technology/log-analysis/internal/interface/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type middlewares struct {
	cfg config.Config
}

func newMiddlewares(cfg config.Config) *middlewares {
	return &middlewares{cfg: cfg}
}

func (m *middlewares) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get token from cookies
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "invalid token")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return c.JSON(http.StatusUnauthorized, "invalid token")
		}
		claims := &models.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.Server.SecretKey), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "invalid token")
		}

		// Store claims in context

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}

// func (m *middlewares) rateLimiterMiddleware() echo.MiddlewareFunc {
// 	cfg := middleware.RateLimiterConfig{
// 		Skipper: middleware.DefaultSkipper,
// 		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
// 			middleware.RateLimiterMemoryStoreConfig{
// 				Rate:      rate.Limit(m.cfg.Server.RateLimiter.Rate),
// 				Burst:     m.cfg.Server.RateLimiter.Burst,
// 				ExpiresIn: m.cfg.Server.RateLimiter.Expires,
// 			}),
// 		IdentifierExtractor: func(c echo.Context) (string, error) {
// 			return c.RealIP(), nil
// 		},
// 		ErrorHandler: func(c echo.Context, err error) error {
// 			slogger.Error(c.Request().Context(), "rate limiter", slog.Any("error", err))
// 			return serde.Response(c, http.StatusInternalServerError, model.InternalMessage, nil)
// 		},
// 		DenyHandler: func(c echo.Context, _ string, _ error) error {
// 			return serde.Response(c, http.StatusTooManyRequests, model.TooManyMessage, nil)
// 		},
// 	}

// 	return middleware.RateLimiterWithConfig(cfg)
// }

func (m *middlewares) corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	})
}

// func (m *middlewares) requestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		id, err := reqid.NewRequestID()
// 		if err != nil {
// 			slogger.Error(c.Request().Context(), "generate request id", slog.Any("error", err))
// 			return serde.Response(c, http.StatusInternalServerError, model.InternalMessage, nil)
// 		}

// 		// put request id inside context
// 		ctx := context.WithValue(c.Request().Context(), reqid.RequestIDKey, id)

// 		// include request_id in logs
// 		ctx = slogger.WithAttrs(ctx, slog.Any("request_id", id))

// 		// replace the context with the new one
// 		c.SetRequest(c.Request().WithContext(ctx))
// 		return next(c)
// 	}
// }

// func (m *middlewares) recoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		defer func() {
// 			if msg := recover(); msg != nil {
// 				slogger.Error(c.Request().Context(), "panic in HTTP server", slog.Any("msg", msg))
// 				c.Response().WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 		}()
// 		return next(c)
// 	}
// }

// func (m *middlewares) loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		start := time.Now()
// 		path := c.Request().URL.Path
// 		raw := c.Request().URL.RawQuery
// 		defer func() {
// 			if raw != "" {
// 				path = path + "?" + raw
// 			}
// 			slogger.Info(c.Request().Context(), "http server",
// 				slog.Group(
// 					"request",
// 					slog.String("client_ip", c.RealIP()),
// 					slog.String("method", c.Request().Method),
// 					slog.String("request_path", path),
// 				),
// 				slog.Group(
// 					"response",
// 					slog.Int("status", c.Response().Status),
// 					slog.String("time_took", time.Since(start).String()),
// 				),
// 			)
// 		}()
// 		return next(c)
// 	}
// }
