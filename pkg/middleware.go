package pkg

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

const CorrelationIDHeaderKey = "X-Correlation-Id"

func RegisterMiddlewares(e *echo.Echo, routePrefix string) {
	e.Pre(AddCorrelationID)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      func(c echo.Context) bool { return shouldSkipMiddleware(c, routePrefix) },
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch,
			http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodTrace, http.MethodConnect,
		},
	}))
	AddHealthCheck(e, routePrefix)
	e.Use(middleware.BodyLimit("1M"))
	e.Use(AddLogger(routePrefix))
	e.Use(AddRecovery(routePrefix))
}

func AddCorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get(CorrelationIDHeaderKey)
		if id == "" {
			id = uuid.New().String()
		}
		c.Request().Header.Set(CorrelationIDHeaderKey, id)
		c.Response().Header().Set(CorrelationIDHeaderKey, id)
		return next(c)
	}
}

func AddLogger(routePrefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shouldSkipMiddleware(c, routePrefix) {
				return next(c)
			}
			return logRequest(c, next)
		}
	}
}

func AddRecovery(routePrefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shouldSkipMiddleware(c, routePrefix) {
				return next(c)
			}
			defer func() {
				if err := recover(); err != nil {
					logrus.WithField("error", err).Error("Panic recovered in middleware")
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"))
				}
			}()
			return next(c)
		}
	}
}

func AddHealthCheck(e *echo.Echo, routePrefix string) {
	e.GET(routePrefix+"/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}

func shouldSkipMiddleware(c echo.Context, routePrefix string) bool {
	path := c.Path()
	return strings.HasPrefix(path, routePrefix+"/swagger") || strings.HasPrefix(path, routePrefix+"/health")
}

func logRequest(c echo.Context, next echo.HandlerFunc) error {
	start := time.Now()
	req := c.Request()
	res := c.Response()

	err := next(c)
	duration := time.Since(start)

	logFields := logrus.Fields{
		"correlationID": req.Header.Get(CorrelationIDHeaderKey),
		"method":        req.Method,
		"path":          req.URL.Path,
		"rawPath":       c.Path(),
		"status":        res.Status,
		"latency":       duration.Microseconds(),
		"latencyHuman":  duration.String(),
		"remoteIP":      c.RealIP(),
		"userAgent":     req.UserAgent(),
		"referer":       req.Referer(),
		"host":          req.Host,
		"error":         err,
	}

	logByStatus(logFields, res.Status)
	return err
}

func logByStatus(fields logrus.Fields, status int) {
	if status >= 500 {
		logrus.WithFields(fields).Error("Request details")
	} else {
		logrus.WithFields(fields).Info("Request details")
	}
}
