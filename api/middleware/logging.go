package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs all incoming requests with request ID, method, path, status, and duration
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate unique request ID for tracing
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		// Add request ID to response headers for client-side tracing
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Record start time
		startTime := time.Now()

		// Log incoming request
		log.Printf("[REQUEST] ID: %s | Method: %s | Path: %s | IP: %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log response with status and duration
		log.Printf("[RESPONSE] ID: %s | Status: %d | Duration: %v | Path: %s",
			requestID,
			c.Writer.Status(),
			duration,
			c.Request.URL.Path,
		)
	}
}
