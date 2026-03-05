package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go-web-demo/pkg/logger"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				requestID := c.GetString(RequestIDKey)

				log.Error("Panic recovered",
					"error", err,
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"stack", string(stack),
				)

				c.AbortWithStatusJSON(500, gin.H{
					"error": fmt.Sprintf("Internal server error: %v", err),
					"request_id": requestID,
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				})
			}
		}()

		c.Next()
	}
}
