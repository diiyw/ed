package handlers

import (
	"encoding/json"
	"os"
	"time"

	"github.com/diiyw/ed/ssh"
)

// SSHConfig represents an SSH server configuration
type SSHConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	AuthType string `json:"auth_type"` // "password", "key", "agent"
	Password string `json:"password,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
	KeyPass  string `json:"key_pass,omitempty"`
}

// Project represents a deployable project
type Project struct {
	Name              string    `json:"name"`
	BuildInstructions string    `json:"build_instructions"`
	DeployScript      string    `json:"deploy_script"`
	DeployServers     []string  `json:"deploy_servers"` // names of SSH configs
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Config holds all application data
type Config struct {
	SSHConfigs []SSHConfig `json:"ssh_configs"`
	Projects   []Project   `json:"projects"`
}

// LoadConfig loads configuration from JSON file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}, nil
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

// SaveConfig saves configuration to JSON file
func SaveConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
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
		return nil, nil
	}
}
