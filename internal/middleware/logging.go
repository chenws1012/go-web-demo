package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go-web-demo/pkg/logger"
)

func Logging(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		entry := log.InfoEvent()
		if statusCode >= 400 {
			entry = log.ErrorEvent()
		}

		if query != "" {
			path = path + "?" + query
		}

		entry.
			Str("request_id", c.GetString(RequestIDKey)).
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("latency", latency.String()).
			Str("user_agent", userAgent).
			Msg("HTTP request")

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error("Request error", "error", e.Error(), "request_id", c.GetString(RequestIDKey))
			}
		}
	}
}

func LoggerWithContext(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString(RequestIDKey)
		subLogger := log.With().Str("request_id", requestID).Logger()

		c.Set("logger", &subLogger)
		c.Next()
	}
}

func GetLogger(c *gin.Context) *zerolog.Logger {
	if logger, exists := c.Get("logger"); exists {
		return logger.(*zerolog.Logger)
	}
	return nil
}
