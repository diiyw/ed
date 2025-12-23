package api

import (
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/diiyw/ed/api/handlers"
)

//go:embed all:testdata/dist
var testEmbeddedFS embed.FS

// TestSetupRouter_AllRoutesRegistered tests that all expected routes are registered
func TestSetupRouter_AllRoutesRegistered(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	router := SetupRouter(config, nil)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		// SSH routes
		{"GET all SSH configs", "GET", "/api/ssh", http.StatusOK},
		{"GET SSH config by name", "GET", "/api/ssh/test", http.StatusNotFound}, // 404 because config doesn't exist
		{"POST create SSH config", "POST", "/api/ssh", http.StatusBadRequest},   // 400 because no body
		{"PUT update SSH config", "PUT", "/api/ssh/test", http.StatusBadRequest},
		{"DELETE SSH config", "DELETE", "/api/ssh/test", http.StatusNotFound},
		{"POST test SSH connection", "POST", "/api/ssh/test/test", http.StatusNotFound},

		// Project routes
		{"GET all projects", "GET", "/api/projects", http.StatusOK},
		{"GET project by name", "GET", "/api/projects/test", http.StatusNotFound},
		{"POST create project", "POST", "/api/projects", http.StatusBadRequest},
		{"PUT update project", "PUT", "/api/projects/test", http.StatusBadRequest},
		{"DELETE project", "DELETE", "/api/projects/test", http.StatusNotFound},
		{"POST deploy project", "POST", "/api/projects/test/deploy", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for %s %s", tt.expectedStatus, w.Code, tt.method, tt.path)
			}
		})
	}
}

// TestSetupRouter_MiddlewareApplied tests that middleware is properly applied
func TestSetupRouter_MiddlewareApplied(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	router := SetupRouter(config, nil)

	t.Run("CORS middleware applied", func(t *testing.T) {
		// Test preflight request
		req, err := http.NewRequest("OPTIONS", "/api/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Origin", "http://localhost:5173")
		req.Header.Set("Access-Control-Request-Method", "GET")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check CORS headers are present
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Error("CORS middleware not applied: Access-Control-Allow-Origin header missing")
		}
	})

	t.Run("Request logger middleware applied", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check that request ID header is present (added by logging middleware)
		if w.Header().Get("X-Request-ID") == "" {
			t.Error("Request logger middleware not applied: X-Request-ID header missing")
		}
	})

	t.Run("Error handler middleware applied", func(t *testing.T) {
		// Error handler is tested by ensuring error responses are properly formatted
		req, err := http.NewRequest("POST", "/api/ssh", strings.NewReader("invalid json"))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get a properly formatted error response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Error response not properly formatted: %v", err)
		}

		if _, ok := response["error"]; !ok {
			t.Error("Error response missing 'error' field")
		}
	})
}

// TestSetupRouter_RouteGroups tests that route groups work correctly
func TestSetupRouter_RouteGroups(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{
			{
				Name:     "test-server",
				Host:     "example.com",
				Port:     22,
				User:     "testuser",
				AuthType: "password",
				Password: "testpass",
			},
		},
		Projects: []handlers.Project{
			{
				Name:              "test-project",
				BuildInstructions: "npm run build",
				DeployScript:      "rsync -avz dist/ server:/var/www/",
				DeployServers:     []string{"test-server"},
			},
		},
	}

	router := SetupRouter(config, nil)

	t.Run("SSH route group", func(t *testing.T) {
		// Test that all SSH routes are under /api/ssh
		req, err := http.NewRequest("GET", "/api/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("SSH route group not working: expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Response missing 'data' field")
		}

		configs, ok := data["configs"].([]interface{})
		if !ok {
			t.Fatal("Response missing 'configs' field")
		}

		if len(configs) != 1 {
			t.Errorf("Expected 1 SSH config, got %d", len(configs))
		}
	})

	t.Run("Projects route group", func(t *testing.T) {
		// Test that all project routes are under /api/projects
		req, err := http.NewRequest("GET", "/api/projects", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Projects route group not working: expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Response missing 'data' field")
		}

		projects, ok := data["projects"].([]interface{})
		if !ok {
			t.Fatal("Response missing 'projects' field")
		}

		if len(projects) != 1 {
			t.Errorf("Expected 1 project, got %d", len(projects))
		}
	})

	t.Run("API base path", func(t *testing.T) {
		// Test that routes without /api prefix don't work
		req, err := http.NewRequest("GET", "/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get 404 because route is under /api/ssh, not /ssh
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404 for route without /api prefix, got %d", w.Code)
		}
	})
}

// TestSetupRouter_Integration tests complete request flow through router
func TestSetupRouter_Integration(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	router := SetupRouter(config, nil)

	t.Run("Complete SSH config workflow", func(t *testing.T) {
		// Create SSH config
		createBody := `{
			"name": "integration-test",
			"host": "example.com",
			"port": 22,
			"user": "testuser",
			"authType": "password",
			"password": "testpass"
		}`

		req, err := http.NewRequest("POST", "/api/ssh", strings.NewReader(createBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Get the created config
		req, err = http.NewRequest("GET", "/api/ssh/integration-test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Complete project workflow", func(t *testing.T) {
		// Create project
		createBody := `{
			"name": "integration-project",
			"buildInstructions": "npm run build",
			"deployScript": "rsync -avz dist/ server:/var/www/",
			"deployServers": []
		}`

		req, err := http.NewRequest("POST", "/api/projects", strings.NewReader(createBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Get the created project
		req, err = http.NewRequest("GET", "/api/projects/integration-project", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

// TestSetupRouter_StaticFileServing tests that embedded static files are served correctly
func TestSetupRouter_StaticFileServing(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	// Create router with embedded test files
	router := SetupRouter(config, &testEmbeddedFS)

	t.Run("Serve index.html at root", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("Response does not contain HTML content")
		}
		if !strings.Contains(body, "Test App") {
			t.Error("Response does not contain expected title")
		}
	})

	t.Run("Serve static JavaScript file", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/assets/main.js", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "Test app loaded") {
			t.Error("JavaScript file content not served correctly")
		}
	})

	t.Run("Serve static CSS file", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/assets/style.css", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "body") {
			t.Error("CSS file content not served correctly")
		}
	})

	t.Run("API routes not affected by static serving", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Should get JSON response, not HTML
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("API route returned non-JSON response: %v", err)
		}
	})

	t.Run("Non-existent API route returns JSON error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/nonexistent", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		// Should get JSON error, not HTML
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Non-existent API route returned non-JSON response: %v", err)
		}

		if _, ok := response["error"]; !ok {
			t.Error("Error response missing 'error' field")
		}
	})
}

// TestSetupRouter_SPAFallback tests that SPA routing fallback works correctly
func TestSetupRouter_SPAFallback(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	router := SetupRouter(config, &testEmbeddedFS)

	t.Run("Non-existent static file falls back to index.html", func(t *testing.T) {
		// Request a route that doesn't exist as a file (SPA route)
		req, err := http.NewRequest("GET", "/projects", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("SPA fallback did not return HTML content")
		}
		if !strings.Contains(body, "Test App") {
			t.Error("SPA fallback did not return index.html")
		}
	})

	t.Run("Nested SPA route falls back to index.html", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/projects/123/edit", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("Nested SPA route did not fall back to index.html")
		}
	})

	t.Run("API routes do not fall back to index.html", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/nonexistent", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get JSON error, not HTML fallback
		body := w.Body.String()
		if strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("API route incorrectly fell back to index.html")
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("API route should return JSON, not HTML: %v", err)
		}
	})
}

// TestSetupRouter_EmbeddedFilesAccessible tests that embedded files are accessible
func TestSetupRouter_EmbeddedFilesAccessible(t *testing.T) {
	config := &handlers.Config{
		SSHConfigs: []handlers.SSHConfig{},
		Projects:   []handlers.Project{},
	}

	t.Run("Router works without embedded files", func(t *testing.T) {
		router := SetupRouter(config, nil)

		req, err := http.NewRequest("GET", "/api/ssh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Router serves embedded files when provided", func(t *testing.T) {
		router := SetupRouter(config, &testEmbeddedFS)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "Test App") {
			t.Error("Embedded files not accessible")
		}
	})

	t.Run("Multiple static files accessible", func(t *testing.T) {
		router := SetupRouter(config, &testEmbeddedFS)

		files := []string{
			"/",
			"/assets/main.js",
			"/assets/style.css",
		}

		for _, file := range files {
			req, err := http.NewRequest("GET", file, nil)
			if err != nil {
				t.Fatalf("Failed to create request for %s: %v", file, err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("File %s not accessible: got status %d", file, w.Code)
			}
		}
	})
}
