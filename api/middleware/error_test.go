package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Unit Tests for Error Middleware
// Requirements: 9.5

func TestPanicRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(ErrorHandler())

	// Handler that panics
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 status code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// Should have error in response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	errorMsg, ok := response["error"].(string)
	if !ok {
		t.Fatal("Response missing 'error' field")
	}

	if !strings.Contains(errorMsg, "test panic") {
		t.Errorf("Error message should contain panic message, got: %s", errorMsg)
	}

	// Should log the panic
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "PANIC RECOVERY") {
		t.Error("Expected panic to be logged with PANIC RECOVERY tag")
	}

	if !strings.Contains(logOutput, "test panic") {
		t.Error("Expected panic message to be logged")
	}

	// Should log context
	if !strings.Contains(logOutput, "PANIC CONTEXT") {
		t.Error("Expected panic context to be logged")
	}
}

func TestErrorResponseFormatting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandler())

	// Handler that returns an error response
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed",
		})
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should preserve the status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// Should have consistent error format
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	errorMsg, ok := response["error"].(string)
	if !ok {
		t.Fatal("Response missing 'error' field")
	}

	if errorMsg != "validation failed" {
		t.Errorf("Expected error message 'validation failed', got: %s", errorMsg)
	}
}

func TestErrorLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(ErrorHandler())

	// Handler that adds an error to the context
	router.GET("/error-context", func(c *gin.Context) {
		c.Error(gin.Error{
			Err:  http.ErrBodyNotAllowed,
			Type: gin.ErrorTypePublic,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	})

	req := httptest.NewRequest("GET", "/error-context", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should log the error with context
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "ERROR") {
		t.Error("Expected error to be logged with ERROR tag")
	}

	if !strings.Contains(logOutput, "GET") {
		t.Error("Expected HTTP method to be logged")
	}

	if !strings.Contains(logOutput, "/error-context") {
		t.Error("Expected request path to be logged")
	}
}

func TestNormalRequestPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandler())

	// Normal handler that doesn't panic or error
	router.GET("/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	req := httptest.NewRequest("GET", "/success", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 200 status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Should have normal response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok {
		t.Fatal("Response missing 'message' field")
	}

	if message != "success" {
		t.Errorf("Expected message 'success', got: %s", message)
	}
}

func TestPanicWithDifferentTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name       string
		panicValue interface{}
	}{
		{"string panic", "string error"},
		{"error panic", http.ErrBodyNotAllowed},
		{"int panic", 42},
		{"struct panic", struct{ msg string }{"error"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			defer log.SetOutput(nil)

			router := gin.New()
			router.Use(ErrorHandler())

			router.GET("/panic", func(c *gin.Context) {
				panic(tc.panicValue)
			})

			req := httptest.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should always return 500
			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500, got %d", w.Code)
			}

			// Should have error in response
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if _, ok := response["error"]; !ok {
				t.Fatal("Response missing 'error' field")
			}

			// Should log the panic
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, "PANIC RECOVERY") {
				t.Error("Expected panic to be logged")
			}
		})
	}
}

func TestMultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(ErrorHandler())

	// Handler that adds multiple errors
	router.GET("/multiple-errors", func(c *gin.Context) {
		c.Error(gin.Error{
			Err:  http.ErrBodyNotAllowed,
			Type: gin.ErrorTypePublic,
		})
		c.Error(gin.Error{
			Err:  http.ErrContentLength,
			Type: gin.ErrorTypePublic,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "multiple errors occurred",
		})
	})

	req := httptest.NewRequest("GET", "/multiple-errors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should log the last error
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "ERROR") {
		t.Error("Expected error to be logged")
	}
}

func TestStackTraceLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	router := gin.New()
	router.Use(ErrorHandler())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic for stack trace")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should log stack trace
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "goroutine") {
		t.Error("Expected stack trace to contain 'goroutine'")
	}
}
