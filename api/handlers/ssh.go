package handlers

import (
	"fmt"
	"net/http"

	"github.com/diiyw/ed/ssh"
	"github.com/gin-gonic/gin"
)

type SSHHandler struct {
	config *Config
}

func NewSSHHandler(config *Config) *SSHHandler {
	return &SSHHandler{config: config}
}

// GetAll returns all SSH configurations
func (h *SSHHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"configs": h.config.SSHConfigs,
		},
	})
}

// GetByName returns a single SSH configuration
func (h *SSHHandler) GetByName(c *gin.Context) {
	name := c.Param("name")

	for _, cfg := range h.config.SSHConfigs {
		if cfg.Name == name {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"config": cfg,
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "SSH configuration not found",
	})
}

// Create creates a new SSH configuration
func (h *SSHHandler) Create(c *gin.Context) {
	var newConfig SSHConfig
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Check if name already exists
	for _, cfg := range h.config.SSHConfigs {
		if cfg.Name == newConfig.Name {
			c.JSON(http.StatusConflict, gin.H{
				"error": "SSH configuration with this name already exists",
			})
			return
		}
	}

	h.config.SSHConfigs = append(h.config.SSHConfigs, newConfig)
	if err := SaveConfig("config.json", h.config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to save configuration: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"config": newConfig,
		},
		"message": "SSH configuration created successfully",
	})
}

// Update updates an existing SSH configuration
func (h *SSHHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var updatedConfig SSHConfig
	if err := c.ShouldBindJSON(&updatedConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	for i, cfg := range h.config.SSHConfigs {
		if cfg.Name == name {
			h.config.SSHConfigs[i] = updatedConfig
			if err := SaveConfig("config.json", h.config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to save configuration: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"config": updatedConfig,
				},
				"message": "SSH configuration updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "SSH configuration not found",
	})
}

// Delete deletes an SSH configuration
func (h *SSHHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	for i, cfg := range h.config.SSHConfigs {
		if cfg.Name == name {
			h.config.SSHConfigs = append(h.config.SSHConfigs[:i], h.config.SSHConfigs[i+1:]...)
			if err := SaveConfig("config.json", h.config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to save configuration: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "SSH configuration deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "SSH configuration not found",
	})
}

// Test tests an SSH connection
func (h *SSHHandler) Test(c *gin.Context) {
	name := c.Param("name")

	var sshConfig *SSHConfig
	for _, cfg := range h.config.SSHConfigs {
		if cfg.Name == name {
			sshConfig = &cfg
			break
		}
	}

	if sshConfig == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "SSH configuration not found",
		})
		return
	}

	// Get auth method
	auth, err := sshConfig.GetAuthMethod()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Auth error: %v", err),
		})
		return
	}

	// Create SSH client
	client, err := ssh.New(sshConfig.User, sshConfig.Host, auth)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("Connection failed: %v", err),
		})
		return
	}
	defer client.Close()

	// Run test command
	output, err := client.Run("echo 'SSH connection successful'")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("Command failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection successful",
		"output":  string(output),
	})
}
