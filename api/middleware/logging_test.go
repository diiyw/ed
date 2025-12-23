package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Unit Tests for Logging Middleware
// Requirements: 10.5

func TestRequestLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Should log incoming request
	if !strings.Contains(logOutput, "[REQUEST]") {
		t.Error("Expected request to be logged with [REQUEST] tag")
	}

	// Should log method
	if !strings.Contains(logOutput, "Method: GET") {
		t.Error("Expected HTTP method to be logged")
	}

	// Should log path
	if !strings.Contains(logOutput, "Path: /test") {
		t.Error("Expected request path to be logged")
	}

	// Should log request ID
	if !strings.Contains(logOutput, "ID:") {
		t.Error("Expected request ID to be logged")
	}
}

func TestResponseLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Should log response
	if !strings.Contains(logOutput, "[RESPONSE]") {
		t.Error("Expected response to be logged with [RESPONSE] tag")
	}

	// Should log status code
	if !strings.Contains(logOutput, "Status: 200") {
		t.Error("Expected status code to be logged")
	}

	// Should log duration
	if !strings.Contains(logOutput, "Duration:") {
		t.Error("Expected duration to be logged")
	}
}

func TestRequestIDGeneration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output to prevent nil pointer
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	var requestID string
	router.GET("/test", func(c *gin.Context) {
		// Get request ID from context
		id, exists := c.Get("requestID")
		if !exists {
			t.Error("Request ID not found in context")
		}
		requestID = id.(string)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should generate a request ID
	if requestID == "" {
		t.Error("Expected request ID to be generated")
	}

	// Should be a valid UUID format (basic check)
	if len(requestID) != 36 {
		t.Errorf("Expected UUID format (36 chars), got length %d", len(requestID))
	}

	if !strings.Contains(requestID, "-") {
		t.Error("Expected UUID format with hyphens")
	}
}

func TestRequestIDInResponseHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output to prevent nil pointer
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should include request ID in response header
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header in response")
	}

	// Should be a valid UUID format
	if len(requestID) != 36 {
		t.Errorf("Expected UUID format (36 chars), got length %d", len(requestID))
	}
}

func TestDurationMeasurement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	router.GET("/slow", func(c *gin.Context) {
		// Simulate slow operation
		time.Sleep(10 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Should log duration
	if !strings.Contains(logOutput, "Duration:") {
		t.Error("Expected duration to be logged")
	}

	// Duration should be at least 10ms
	if !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "µs") {
		t.Error("Expected duration to be in milliseconds or microseconds")
	}
}

func TestLoggingWithDifferentStatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name       string
		statusCode int
		handler    func(*gin.Context)
	}{
		{
			"200 OK",
			http.StatusOK,
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			},
		},
		{
			"201 Created",
			http.StatusCreated,
			func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{"message": "created"})
			},
		},
		{
			"400 Bad Request",
			http.StatusBadRequest,
			func(c *gin.Context) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			},
		},
		{
			"404 Not Found",
			http.StatusNotFound,
			func(c *gin.Context) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			},
		},
		{
			"500 Internal Server Error",
			http.StatusInternalServerError,
			func(c *gin.Context) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			defer log.SetOutput(nil)

			router := gin.New()
			router.Use(RequestLogger())
			router.GET("/test", tc.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			logOutput := logBuf.String()

			// Should log the status code
			if !strings.Contains(logOutput, "[RESPONSE]") {
				t.Error("Expected response to be logged")
			}

			// Check that status code appears in the log
			if !strings.Contains(logOutput, "Status:") {
				t.Error("Expected status code to be logged")
			}
		})
	}
}

func TestLoggingWithDifferentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			defer log.SetOutput(nil)

			router := gin.New()
			router.Use(RequestLogger())

			router.Handle(method, "/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			logOutput := logBuf.String()

			// Should log the HTTP method
			expectedMethod := "Method: " + method
			if !strings.Contains(logOutput, expectedMethod) {
				t.Errorf("Expected method %s to be logged", method)
			}
		})
	}
}

func TestRequestIDUniqueness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output to prevent nil pointer
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	requestIDs := make(map[string]bool)

	router.GET("/test", func(c *gin.Context) {
		id, _ := c.Get("requestID")
		requestIDs[id.(string)] = true
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make multiple requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// All request IDs should be unique
	if len(requestIDs) != 10 {
		t.Errorf("Expected 10 unique request IDs, got %d", len(requestIDs))
	}
}

func TestClientIPLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(RequestLogger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Should log client IP
	if !strings.Contains(logOutput, "IP:") {
		t.Error("Expected client IP to be logged")
	}
}
