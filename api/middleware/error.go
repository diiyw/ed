package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorHandler provides panic recovery and consistent error response formatting
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				stack := debug.Stack()
				log.Printf("[PANIC RECOVERY] %v\n%s", err, stack)

				// Log additional context
				log.Printf("[PANIC CONTEXT] Method: %s, Path: %s, IP: %s",
					c.Request.Method,
					c.Request.URL.Path,
					c.ClientIP(),
				)

				// Format error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Internal server error: %v", err),
				})

				// Abort further processing
				c.Abort()
			}
		}()

		// Process request
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			// Log error with context
			log.Printf("[ERROR] Method: %s, Path: %s, IP: %s, Error: %v",
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				err.Err,
			)

			// If response hasn't been written yet, format error response
			if !c.Writer.Written() {
				statusCode := c.Writer.Status()
				if statusCode == http.StatusOK {
					statusCode = http.StatusInternalServerError
				}

				c.JSON(statusCode, gin.H{
					"error": err.Error(),
				})
			}
		}
	}
}
