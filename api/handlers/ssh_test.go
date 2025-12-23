package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: frontend-refactor, Property 18: HTTP status codes for errors
// Validates: Requirements 9.5
func TestProperty_HTTPStatusCodesForErrors(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	// Property: For any API endpoint that encounters an error, the response should include
	// an appropriate HTTP status code (4xx for client errors, 5xx for server errors)
	properties.Property("API endpoints return appropriate error status codes", prop.ForAll(
		func(names []string) bool {
			gin.SetMode(gin.TestMode)

			// Test various error scenarios
			testCases := []struct {
				name           string
				setupConfig    func() *Config
				method         string
				path           string
				body           interface{}
				expectedStatus int
				statusRange    string // "4xx" or "5xx"
			}{
				{
					name: "GET non-existent SSH config returns 404",
					setupConfig: func() *Config {
						return &Config{SSHConfigs: []SSHConfig{}}
					},
					method:         "GET",
					path:           "/api/ssh/nonexistent",
					expectedStatus: http.StatusNotFound,
					statusRange:    "4xx",
				},
				{
					name: "POST invalid JSON returns 400",
					setupConfig: func() *Config {
						return &Config{SSHConfigs: []SSHConfig{}}
					},
					method:         "POST",
					path:           "/api/ssh",
					body:           "invalid json",
					expectedStatus: http.StatusBadRequest,
					statusRange:    "4xx",
				},
				{
					name: "POST duplicate name returns 409",
					setupConfig: func() *Config {
						return &Config{
							SSHConfigs: []SSHConfig{
								{Name: names[0], Host: "test.com", Port: 22, User: "user", AuthType: "password"},
							},
						}
					},
					method: "POST",
					path:   "/api/ssh",
					body: SSHConfig{
						Name:     names[0],
						Host:     "test2.com",
						Port:     22,
						User:     "user",
						AuthType: "password",
					},
					expectedStatus: http.StatusConflict,
					statusRange:    "4xx",
				},
				{
					name: "PUT non-existent config returns 404",
					setupConfig: func() *Config {
						return &Config{SSHConfigs: []SSHConfig{}}
					},
					method: "PUT",
					path:   "/api/ssh/nonexistent",
					body: SSHConfig{
						Name:     "nonexistent",
						Host:     "test.com",
						Port:     22,
						User:     "user",
						AuthType: "password",
					},
					expectedStatus: http.StatusNotFound,
					statusRange:    "4xx",
				},
				{
					name: "DELETE non-existent config returns 404",
					setupConfig: func() *Config {
						return &Config{SSHConfigs: []SSHConfig{}}
					},
					method:         "DELETE",
					path:           "/api/ssh/nonexistent",
					expectedStatus: http.StatusNotFound,
					statusRange:    "4xx",
				},
				{
					name: "POST test non-existent config returns 404",
					setupConfig: func() *Config {
						return &Config{SSHConfigs: []SSHConfig{}}
					},
					method:         "POST",
					path:           "/api/ssh/nonexistent/test",
					expectedStatus: http.StatusNotFound,
					statusRange:    "4xx",
				},
			}

			for _, tc := range testCases {
				config := tc.setupConfig()
				handler := NewSSHHandler(config)

				router := gin.New()
				router.GET("/api/ssh/:name", handler.GetByName)
				router.POST("/api/ssh", handler.Create)
				router.POST("/api/ssh/:name/test", handler.Test)
				router.PUT("/api/ssh/:name", handler.Update)
				router.DELETE("/api/ssh/:name", handler.Delete)

				var req *http.Request
				if tc.body != nil {
					if str, ok := tc.body.(string); ok {
						req = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(str))
					} else {
						bodyBytes, _ := json.Marshal(tc.body)
						req = httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyBytes))
					}
					req.Header.Set("Content-Type", "application/json")
				} else {
					req = httptest.NewRequest(tc.method, tc.path, nil)
				}

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Check that status code is in the expected range
				statusCode := w.Code
				if tc.statusRange == "4xx" {
					if statusCode < 400 || statusCode >= 500 {
						t.Logf("Test case '%s' failed: expected 4xx status code, got %d", tc.name, statusCode)
						return false
					}
				} else if tc.statusRange == "5xx" {
					if statusCode < 500 || statusCode >= 600 {
						t.Logf("Test case '%s' failed: expected 5xx status code, got %d", tc.name, statusCode)
						return false
					}
				}

				// Verify the response has an error field
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Logf("Test case '%s' failed: could not parse response JSON", tc.name)
					return false
				}

				if _, hasError := response["error"]; !hasError {
					t.Logf("Test case '%s' failed: response missing 'error' field", tc.name)
					return false
				}
			}

			return true
		},
		gen.SliceOfN(5, gen.Identifier()).SuchThat(func(v interface{}) bool {
			// Generate non-empty names
			names := v.([]string)
			return len(names) > 0
		}),
	))

	properties.TestingRun(t)
}

// Unit Tests for SSH API Handlers
// Requirements: 9.1, 9.3

func TestGetAllSSHConfigs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
			{Name: "server2", Host: "host2.com", Port: 22, User: "user2", AuthType: "key"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.GET("/api/ssh", handler.GetAll)

	req := httptest.NewRequest("GET", "/api/ssh", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
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

	if len(configs) != 2 {
		t.Errorf("Expected 2 configs, got %d", len(configs))
	}
}

func TestCreateSSHConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.POST("/api/ssh", handler.Create)

	newConfig := SSHConfig{
		Name:     "newserver",
		Host:     "newhost.com",
		Port:     22,
		User:     "newuser",
		AuthType: "password",
		Password: "secret",
	}

	bodyBytes, _ := json.Marshal(newConfig)
	req := httptest.NewRequest("POST", "/api/ssh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response missing 'message' field")
	}

	if _, ok := response["data"]; !ok {
		t.Error("Response missing 'data' field")
	}
}

func TestUpdateSSHConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.PUT("/api/ssh/:name", handler.Update)

	updatedConfig := SSHConfig{
		Name:     "server1",
		Host:     "updated.com",
		Port:     2222,
		User:     "updateduser",
		AuthType: "key",
		KeyFile:  "/path/to/key",
	}

	bodyBytes, _ := json.Marshal(updatedConfig)
	req := httptest.NewRequest("PUT", "/api/ssh/server1", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response missing 'message' field")
	}

	// Verify the config was actually updated
	if config.SSHConfigs[0].Host != "updated.com" {
		t.Errorf("Config was not updated, host is still %s", config.SSHConfigs[0].Host)
	}
}

func TestDeleteSSHConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
			{Name: "server2", Host: "host2.com", Port: 22, User: "user2", AuthType: "key"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.DELETE("/api/ssh/:name", handler.Delete)

	req := httptest.NewRequest("DELETE", "/api/ssh/server1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["message"]; !ok {
		t.Error("Response missing 'message' field")
	}

	// Verify the config was actually deleted
	if len(config.SSHConfigs) != 1 {
		t.Errorf("Expected 1 config remaining, got %d", len(config.SSHConfigs))
	}

	if config.SSHConfigs[0].Name == "server1" {
		t.Error("server1 was not deleted")
	}
}

func TestTestSSHConnection_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.POST("/api/ssh/:name/test", handler.Test)

	req := httptest.NewRequest("POST", "/api/ssh/nonexistent/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response missing 'error' field")
	}
}

func TestTestSSHConnection_InvalidAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use a config with key auth but invalid key file path
	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "key", KeyFile: "/nonexistent/key"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.POST("/api/ssh/:name/test", handler.Test)

	req := httptest.NewRequest("POST", "/api/ssh/server1/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 400 for invalid key file
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	success, ok := response["success"].(bool)
	if !ok || success {
		t.Error("Expected success to be false")
	}
}

func TestGetSSHConfigByName_Found(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.GET("/api/ssh/:name", handler.GetByName)

	req := httptest.NewRequest("GET", "/api/ssh/server1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
	}

	configData, ok := data["config"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'config' field")
	}

	if configData["name"] != "server1" {
		t.Errorf("Expected name 'server1', got %v", configData["name"])
	}
}

func TestGetSSHConfigByName_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.GET("/api/ssh/:name", handler.GetByName)

	req := httptest.NewRequest("GET", "/api/ssh/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response missing 'error' field")
	}
}

func TestCreateSSHConfig_DuplicateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.POST("/api/ssh", handler.Create)

	duplicateConfig := SSHConfig{
		Name:     "server1",
		Host:     "different.com",
		Port:     22,
		User:     "user",
		AuthType: "password",
	}

	bodyBytes, _ := json.Marshal(duplicateConfig)
	req := httptest.NewRequest("POST", "/api/ssh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response missing 'error' field")
	}
}

func TestCreateSSHConfig_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.POST("/api/ssh", handler.Create)

	req := httptest.NewRequest("POST", "/api/ssh", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response missing 'error' field")
	}
}

func TestUpdateSSHConfig_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
		},
	}

	handler := NewSSHHandler(config)
	router := gin.New()
	router.PUT("/api/ssh/:name", handler.Update)

	req := httptest.NewRequest("PUT", "/api/ssh/server1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("Response missing 'error' field")
	}
}
