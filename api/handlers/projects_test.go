package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Unit Tests for Project API Handlers
// Requirements: 9.2, 9.4

func TestGetAllProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install && npm run build",
				DeployScript:      "rsync -avz dist/ server:/var/www/",
				DeployServers:     []string{"server1"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				Name:              "project2",
				BuildInstructions: "go build",
				DeployScript:      "scp binary server:/usr/local/bin/",
				DeployServers:     []string{"server2"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.GET("/api/projects", handler.GetAll)

	req := httptest.NewRequest("GET", "/api/projects", nil)
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

	projects, ok := data["projects"].([]interface{})
	if !ok {
		t.Fatal("Response missing 'projects' field")
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestCreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects", handler.Create)

	newProject := Project{
		Name:              "newproject",
		BuildInstructions: "npm install && npm run build",
		DeployScript:      "rsync -avz dist/ server:/var/www/",
		DeployServers:     []string{"server1", "server2"},
	}

	bodyBytes, _ := json.Marshal(newProject)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(bodyBytes))
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

	// Verify the project was added to config
	if len(config.Projects) != 1 {
		t.Errorf("Expected 1 project in config, got %d", len(config.Projects))
	}

	// Verify timestamps were set
	if config.Projects[0].CreatedAt.IsZero() {
		t.Error("CreatedAt timestamp was not set")
	}
	if config.Projects[0].UpdatedAt.IsZero() {
		t.Error("UpdatedAt timestamp was not set")
	}
}

func TestUpdateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt := time.Now().Add(-24 * time.Hour)
	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1"},
				CreatedAt:         createdAt,
				UpdatedAt:         createdAt,
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.PUT("/api/projects/:name", handler.Update)

	updatedProject := Project{
		Name:              "project1",
		BuildInstructions: "npm install && npm run build",
		DeployScript:      "rsync -avz dist/ server:/var/www/",
		DeployServers:     []string{"server1", "server2"},
	}

	bodyBytes, _ := json.Marshal(updatedProject)
	req := httptest.NewRequest("PUT", "/api/projects/project1", bytes.NewBuffer(bodyBytes))
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

	// Verify the project was actually updated
	if config.Projects[0].BuildInstructions != "npm install && npm run build" {
		t.Errorf("Project was not updated, build instructions is still %s", config.Projects[0].BuildInstructions)
	}

	// Verify CreatedAt was preserved
	if !config.Projects[0].CreatedAt.Equal(createdAt) {
		t.Error("CreatedAt timestamp was modified")
	}

	// Verify UpdatedAt was changed
	if config.Projects[0].UpdatedAt.Equal(createdAt) {
		t.Error("UpdatedAt timestamp was not updated")
	}
}

func TestDeleteProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				Name:              "project2",
				BuildInstructions: "go build",
				DeployScript:      "scp binary",
				DeployServers:     []string{"server2"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.DELETE("/api/projects/:name", handler.Delete)

	req := httptest.NewRequest("DELETE", "/api/projects/project1", nil)
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

	// Verify the project was actually deleted
	if len(config.Projects) != 1 {
		t.Errorf("Expected 1 project remaining, got %d", len(config.Projects))
	}

	if config.Projects[0].Name == "project1" {
		t.Error("project1 was not deleted")
	}
}

func TestDeployProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
			{Name: "server2", Host: "host2.com", Port: 22, User: "user2", AuthType: "key"},
		},
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install && npm run build",
				DeployScript:      "rsync -avz dist/ server:/var/www/",
				DeployServers:     []string{"server1", "server2"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects/:name/deploy", handler.Deploy)

	req := httptest.NewRequest("POST", "/api/projects/project1/deploy", nil)
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

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'data' field")
	}

	deploymentID, ok := data["deploymentId"].(string)
	if !ok || deploymentID == "" {
		t.Error("Response missing or empty 'deploymentId' field")
	}
}

func TestGetProjectByName_Found(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.GET("/api/projects/:name", handler.GetByName)

	req := httptest.NewRequest("GET", "/api/projects/project1", nil)
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

	projectData, ok := data["project"].(map[string]interface{})
	if !ok {
		t.Fatal("Response missing 'project' field")
	}

	if projectData["name"] != "project1" {
		t.Errorf("Expected name 'project1', got %v", projectData["name"])
	}
}

func TestGetProjectByName_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.GET("/api/projects/:name", handler.GetByName)

	req := httptest.NewRequest("GET", "/api/projects/nonexistent", nil)
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

func TestCreateProject_DuplicateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects", handler.Create)

	duplicateProject := Project{
		Name:              "project1",
		BuildInstructions: "different build",
		DeployScript:      "different deploy",
		DeployServers:     []string{"server2"},
	}

	bodyBytes, _ := json.Marshal(duplicateProject)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(bodyBytes))
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

func TestCreateProject_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects", handler.Create)

	req := httptest.NewRequest("POST", "/api/projects", bytes.NewBufferString("invalid json"))
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

func TestUpdateProject_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.PUT("/api/projects/:name", handler.Update)

	req := httptest.NewRequest("PUT", "/api/projects/project1", bytes.NewBufferString("invalid json"))
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

func TestUpdateProject_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.PUT("/api/projects/:name", handler.Update)

	updatedProject := Project{
		Name:              "nonexistent",
		BuildInstructions: "npm install",
		DeployScript:      "rsync dist/",
		DeployServers:     []string{"server1"},
	}

	bodyBytes, _ := json.Marshal(updatedProject)
	req := httptest.NewRequest("PUT", "/api/projects/nonexistent", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
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

func TestDeleteProject_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.DELETE("/api/projects/:name", handler.Delete)

	req := httptest.NewRequest("DELETE", "/api/projects/nonexistent", nil)
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

func TestDeployProject_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		Projects: []Project{},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects/:name/deploy", handler.Deploy)

	req := httptest.NewRequest("POST", "/api/projects/nonexistent/deploy", nil)
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

func TestDeployProject_InvalidServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password"},
		},
		Projects: []Project{
			{
				Name:              "project1",
				BuildInstructions: "npm install",
				DeployScript:      "rsync dist/",
				DeployServers:     []string{"server1", "nonexistent"},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		},
	}

	handler := NewProjectHandler(config)
	router := gin.New()
	router.POST("/api/projects/:name/deploy", handler.Deploy)

	req := httptest.NewRequest("POST", "/api/projects/project1/deploy", nil)
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
