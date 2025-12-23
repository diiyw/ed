package deploy

import (
	"context"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}
	if engine.deployments == nil {
		t.Error("deployments map not initialized")
	}
	if engine.clients == nil {
		t.Error("clients map not initialized")
	}
}

func TestGetStatus(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Test non-existent deployment
	_, exists := engine.GetStatus(deploymentID)
	if exists {
		t.Error("Expected deployment to not exist")
	}

	// Add a deployment
	engine.deployments[deploymentID] = &DeploymentStatus{
		ID:          deploymentID,
		ProjectName: "test-project",
		Status:      "running",
		StartedAt:   time.Now(),
	}

	// Test existing deployment
	status, exists := engine.GetStatus(deploymentID)
	if !exists {
		t.Error("Expected deployment to exist")
	}
	if status.ID != deploymentID {
		t.Errorf("Expected ID %s, got %s", deploymentID, status.ID)
	}
	if status.Status != "running" {
		t.Errorf("Expected status 'running', got %s", status.Status)
	}
}

func TestRegisterAndUnregisterClient(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Create mock WebSocket connection (nil is acceptable for this test)
	var conn1, conn2 *websocket.Conn

	// Register first client
	engine.RegisterClient(deploymentID, conn1)
	if len(engine.clients[deploymentID]) != 1 {
		t.Errorf("Expected 1 client, got %d", len(engine.clients[deploymentID]))
	}

	// Register second client
	engine.RegisterClient(deploymentID, conn2)
	if len(engine.clients[deploymentID]) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(engine.clients[deploymentID]))
	}

	// Unregister first client
	engine.UnregisterClient(deploymentID, conn1)
	if len(engine.clients[deploymentID]) != 1 {
		t.Errorf("Expected 1 client after unregister, got %d", len(engine.clients[deploymentID]))
	}

	// Unregister second client
	engine.UnregisterClient(deploymentID, conn2)
	if _, exists := engine.clients[deploymentID]; exists {
		t.Error("Expected clients map entry to be deleted when no clients remain")
	}
}

func TestExecuteBuild(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"
	ctx := context.Background()

	buildInstructions := `
# This is a comment
npm install
npm run build
npm test
`

	err := engine.executeBuild(ctx, deploymentID, buildInstructions)
	if err != nil {
		t.Errorf("executeBuild failed: %v", err)
	}
}

func TestExecuteBuildWithCancellation(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"
	ctx, cancel := context.WithCancel(context.Background())

	buildInstructions := `
npm install
npm run build
npm test
`

	// Cancel immediately
	cancel()

	err := engine.executeBuild(ctx, deploymentID, buildInstructions)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestCompleteDeployment(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Initialize deployment
	engine.deployments[deploymentID] = &DeploymentStatus{
		ID:          deploymentID,
		ProjectName: "test-project",
		Status:      "running",
		StartedAt:   time.Now(),
	}

	// Complete deployment
	engine.completeDeployment(deploymentID)

	status, exists := engine.GetStatus(deploymentID)
	if !exists {
		t.Fatal("Deployment not found")
	}
	if status.Status != "success" {
		t.Errorf("Expected status 'success', got %s", status.Status)
	}
	if status.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestFailDeployment(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Initialize deployment
	engine.deployments[deploymentID] = &DeploymentStatus{
		ID:          deploymentID,
		ProjectName: "test-project",
		Status:      "running",
		StartedAt:   time.Now(),
	}

	// Fail deployment
	errorMsg := "Build failed: command not found"
	engine.failDeployment(deploymentID, errorMsg)

	status, exists := engine.GetStatus(deploymentID)
	if !exists {
		t.Fatal("Deployment not found")
	}
	if status.Status != "failed" {
		t.Errorf("Expected status 'failed', got %s", status.Status)
	}
	if status.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestCancelDeployment(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Test cancelling non-existent deployment
	err := engine.CancelDeployment(deploymentID)
	if err == nil {
		t.Error("Expected error when cancelling non-existent deployment")
	}

	// Initialize deployment
	engine.deployments[deploymentID] = &DeploymentStatus{
		ID:          deploymentID,
		ProjectName: "test-project",
		Status:      "running",
		StartedAt:   time.Now(),
	}

	// Cancel deployment
	err = engine.CancelDeployment(deploymentID)
	if err != nil {
		t.Errorf("CancelDeployment failed: %v", err)
	}

	status, exists := engine.GetStatus(deploymentID)
	if !exists {
		t.Fatal("Deployment not found")
	}
	if status.Status != "failed" {
		t.Errorf("Expected status 'failed', got %s", status.Status)
	}
	if status.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}

	// Test cancelling already completed deployment
	err = engine.CancelDeployment(deploymentID)
	if err == nil {
		t.Error("Expected error when cancelling non-running deployment")
	}
}

func TestGetAuthMethod(t *testing.T) {
	tests := []struct {
		name      string
		config    SSHConfig
		expectErr bool
	}{
		{
			name: "password auth",
			config: SSHConfig{
				AuthType: "password",
				Password: "test123",
			},
			expectErr: false,
		},
		{
			name: "unknown auth type",
			config: SSHConfig{
				AuthType: "unknown",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := tt.config.GetAuthMethod()
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if auth == nil {
					t.Error("Expected auth method but got nil")
				}
			}
		})
	}
}

func TestDeploymentStatusTracking(t *testing.T) {
	engine := NewEngine()
	deploymentID := "test-deployment-1"

	// Create a simple project
	project := &Project{
		Name:              "test-project",
		BuildInstructions: "echo 'Building...'",
		DeployScript:      "",
		DeployServers:     []string{},
	}

	sshConfigs := make(map[string]*SSHConfig)

	// Start deployment in a goroutine
	ctx := context.Background()
	go func() {
		_ = engine.Deploy(ctx, deploymentID, project, sshConfigs)
	}()

	// Wait a bit for deployment to initialize
	time.Sleep(100 * time.Millisecond)

	// Check that deployment was created
	status, exists := engine.GetStatus(deploymentID)
	if !exists {
		t.Fatal("Deployment not found")
	}
	if status.ProjectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got %s", status.ProjectName)
	}

	// Wait for deployment to complete
	time.Sleep(500 * time.Millisecond)

	// Check final status
	status, exists = engine.GetStatus(deploymentID)
	if !exists {
		t.Fatal("Deployment not found")
	}
	if status.Status != "success" {
		t.Errorf("Expected status 'success', got %s", status.Status)
	}
}
