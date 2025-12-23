package deploy

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/diiyw/ed/ssh"
	"github.com/gorilla/websocket"
	xssh "golang.org/x/crypto/ssh"
)

// LogType represents the type of log message
type LogType string

const (
	LogTypeLog    LogType = "log"
	LogTypeStatus LogType = "status"
	LogTypeError  LogType = "error"
)

// DeploymentLog represents a log entry from deployment
type DeploymentLog struct {
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// DeploymentStatus represents the status of a deployment
type DeploymentStatus struct {
	ID          string     `json:"id"`
	ProjectName string     `json:"projectName"`
	Status      string     `json:"status"` // "pending", "running", "success", "failed"
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// SSHConfig represents an SSH server configuration
type SSHConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	AuthType string `json:"auth_type"`
	Password string `json:"password,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
	KeyPass  string `json:"key_pass,omitempty"`
}

// Project represents a deployable project
type Project struct {
	Name              string    `json:"name"`
	BuildInstructions string    `json:"build_instructions"`
	DeployScript      string    `json:"deploy_script"`
	DeployServers     []string  `json:"deploy_servers"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Engine manages deployment operations
type Engine struct {
	deployments map[string]*DeploymentStatus
	clients     map[string][]*websocket.Conn
	mu          sync.RWMutex
}

// NewEngine creates a new deployment engine
func NewEngine() *Engine {
	return &Engine{
		deployments: make(map[string]*DeploymentStatus),
		clients:     make(map[string][]*websocket.Conn),
	}
}

// GetAuthMethod returns the SSH auth method for this config
func (sc *SSHConfig) GetAuthMethod() (ssh.Auth, error) {
	switch sc.AuthType {
	case "password":
		return ssh.Password(sc.Password), nil
	case "key":
		return ssh.Key(sc.KeyFile, sc.KeyPass)
	case "agent":
		return ssh.UseAgent()
	default:
		return nil, fmt.Errorf("unknown auth type: %s", sc.AuthType)
	}
}

// RegisterClient registers a WebSocket client for a deployment
func (e *Engine) RegisterClient(deploymentID string, conn *websocket.Conn) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.clients[deploymentID] = append(e.clients[deploymentID], conn)
}

// UnregisterClient removes a WebSocket client
func (e *Engine) UnregisterClient(deploymentID string, conn *websocket.Conn) {
	e.mu.Lock()
	defer e.mu.Unlock()

	clients := e.clients[deploymentID]
	for i, c := range clients {
		if c == conn {
			e.clients[deploymentID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(e.clients[deploymentID]) == 0 {
		delete(e.clients, deploymentID)
	}
}

// broadcastLog sends a log message to all connected clients
func (e *Engine) broadcastLog(deploymentID string, logType LogType, message string) {
	e.mu.RLock()
	clients := e.clients[deploymentID]
	e.mu.RUnlock()

	log := DeploymentLog{
		Type:      string(logType),
		Data:      message,
		Timestamp: time.Now(),
	}

	for _, conn := range clients {
		if err := conn.WriteJSON(log); err != nil {
			// Client disconnected, will be cleaned up later
			continue
		}
	}
}

// GetStatus returns the status of a deployment
func (e *Engine) GetStatus(deploymentID string) (*DeploymentStatus, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	status, exists := e.deployments[deploymentID]
	return status, exists
}

// Deploy executes a deployment
func (e *Engine) Deploy(ctx context.Context, deploymentID string, project *Project, sshConfigs map[string]*SSHConfig) error {
	// Initialize deployment status
	e.mu.Lock()
	e.deployments[deploymentID] = &DeploymentStatus{
		ID:          deploymentID,
		ProjectName: project.Name,
		Status:      "running",
		StartedAt:   time.Now(),
	}
	e.mu.Unlock()

	e.broadcastLog(deploymentID, LogTypeStatus, "Deployment started")
	e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("Project: %s", project.Name))

	// Execute build instructions if provided
	if project.BuildInstructions != "" {
		e.broadcastLog(deploymentID, LogTypeLog, "Executing build instructions...")
		if err := e.executeBuild(ctx, deploymentID, project.BuildInstructions); err != nil {
			e.failDeployment(deploymentID, fmt.Sprintf("Build failed: %v", err))
			return err
		}
		e.broadcastLog(deploymentID, LogTypeLog, "Build completed successfully")
	}

	// Deploy to each server
	for _, serverName := range project.DeployServers {
		sshConfig, exists := sshConfigs[serverName]
		if !exists {
			err := fmt.Errorf("SSH config not found: %s", serverName)
			e.failDeployment(deploymentID, err.Error())
			return err
		}

		e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("Deploying to server: %s", serverName))
		if err := e.deployToServer(ctx, deploymentID, project, sshConfig); err != nil {
			e.failDeployment(deploymentID, fmt.Sprintf("Deployment to %s failed: %v", serverName, err))
			return err
		}
		e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("Successfully deployed to %s", serverName))
	}

	// Mark deployment as successful
	e.completeDeployment(deploymentID)
	return nil
}

// executeBuild executes build instructions locally
func (e *Engine) executeBuild(ctx context.Context, deploymentID string, buildInstructions string) error {
	// For now, we'll just log the build instructions
	// In a real implementation, this would execute the build commands
	lines := strings.Split(buildInstructions, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("  > %s", line))

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Simulate build time
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// deployToServer deploys to a single SSH server
func (e *Engine) deployToServer(ctx context.Context, deploymentID string, project *Project, sshConfig *SSHConfig) error {
	// Get SSH auth method
	auth, err := sshConfig.GetAuthMethod()
	if err != nil {
		return fmt.Errorf("failed to get auth method: %w", err)
	}

	// Create SSH client
	client, err := ssh.NewConn(&ssh.Config{
		User:     sshConfig.User,
		Addr:     sshConfig.Host,
		Port:     uint(sshConfig.Port),
		Auth:     auth,
		Timeout:  30 * time.Second,
		Callback: xssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer client.Close()

	e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("Connected to %s@%s:%d", sshConfig.User, sshConfig.Host, sshConfig.Port))

	// Execute deploy script
	if project.DeployScript != "" {
		lines := strings.Split(project.DeployScript, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Check for context cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			e.broadcastLog(deploymentID, LogTypeLog, fmt.Sprintf("  $ %s", line))

			// Execute command
			cmd, err := client.CommandContext(ctx, "bash", "-c", line)
			if err != nil {
				return fmt.Errorf("failed to create command: %w", err)
			}

			// Capture output
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("failed to get stdout pipe: %w", err)
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				return fmt.Errorf("failed to get stderr pipe: %w", err)
			}

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start command: %w", err)
			}

			// Stream output
			go e.streamOutput(deploymentID, stdout, LogTypeLog)
			go e.streamOutput(deploymentID, stderr, LogTypeError)

			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("command failed: %w", err)
			}
		}
	}

	return nil
}

// streamOutput streams command output to WebSocket clients
func (e *Engine) streamOutput(deploymentID string, reader io.Reader, logType LogType) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			output := strings.TrimSpace(string(buf[:n]))
			if output != "" {
				e.broadcastLog(deploymentID, logType, output)
			}
		}
		if err != nil {
			if err != io.EOF {
				e.broadcastLog(deploymentID, LogTypeError, fmt.Sprintf("Error reading output: %v", err))
			}
			break
		}
	}
}

// completeDeployment marks a deployment as successful
func (e *Engine) completeDeployment(deploymentID string) {
	e.mu.Lock()
	if status, exists := e.deployments[deploymentID]; exists {
		now := time.Now()
		status.Status = "success"
		status.CompletedAt = &now
	}
	e.mu.Unlock()

	e.broadcastLog(deploymentID, LogTypeStatus, "Deployment completed successfully")
}

// failDeployment marks a deployment as failed
func (e *Engine) failDeployment(deploymentID string, errorMsg string) {
	e.mu.Lock()
	if status, exists := e.deployments[deploymentID]; exists {
		now := time.Now()
		status.Status = "failed"
		status.CompletedAt = &now
	}
	e.mu.Unlock()

	e.broadcastLog(deploymentID, LogTypeError, errorMsg)
	e.broadcastLog(deploymentID, LogTypeStatus, "Deployment failed")
}

// CancelDeployment cancels an ongoing deployment
func (e *Engine) CancelDeployment(deploymentID string) error {
	e.mu.Lock()
	status, exists := e.deployments[deploymentID]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("deployment not found")
	}

	if status.Status != "running" {
		e.mu.Unlock()
		return fmt.Errorf("deployment is not running")
	}

	now := time.Now()
	status.Status = "failed"
	status.CompletedAt = &now
	e.mu.Unlock()

	e.broadcastLog(deploymentID, LogTypeStatus, "Deployment cancelled by user")
	return nil
}
