package middleware

import (
	logger "testingtask/pkg"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := uuid.New().String()

		l := logger.With(map[string]interface{}{
			"request_id": reqID,
			"method":     c.Request().Method,
			"path":       c.Request().URL.Path,
			"remote_ip":  c.RealIP(),
		})

		ctx := logger.WithContext(c.Request().Context(), l)
		c.SetRequest(c.Request().WithContext(ctx))

		start := time.Now()
		err := next(c)
		duration := time.Since(start)

		logger.Info(ctx, "request completed", map[string]interface{}{
			"status":      c.Response().Status,
			"duration_ms": duration.Milliseconds(),
		})

		return err
	}
}
