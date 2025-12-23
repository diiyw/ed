package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	config *Config
}

func NewProjectHandler(config *Config) *ProjectHandler {
	return &ProjectHandler{config: config}
}

// GetAll returns all projects
func (h *ProjectHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"projects": h.config.Projects,
		},
	})
}

// GetByName returns a single project
func (h *ProjectHandler) GetByName(c *gin.Context) {
	name := c.Param("name")

	for _, proj := range h.config.Projects {
		if proj.Name == name {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"project": proj,
				},
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Project not found",
	})
}

// Create creates a new project
func (h *ProjectHandler) Create(c *gin.Context) {
	var newProject Project
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Check if name already exists
	for _, proj := range h.config.Projects {
		if proj.Name == newProject.Name {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Project with this name already exists",
			})
			return
		}
	}

	// Set timestamps
	now := time.Now()
	newProject.CreatedAt = now
	newProject.UpdatedAt = now

	h.config.Projects = append(h.config.Projects, newProject)
	if err := SaveConfig("config.json", h.config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to save configuration: %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"project": newProject,
		},
		"message": "Project created successfully",
	})
}

// Update updates an existing project
func (h *ProjectHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var updatedProject Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	for i, proj := range h.config.Projects {
		if proj.Name == name {
			// Preserve CreatedAt, update UpdatedAt
			updatedProject.CreatedAt = proj.CreatedAt
			updatedProject.UpdatedAt = time.Now()

			h.config.Projects[i] = updatedProject
			if err := SaveConfig("config.json", h.config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to save configuration: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"project": updatedProject,
				},
				"message": "Project updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Project not found",
	})
}

// Delete deletes a project
func (h *ProjectHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	for i, proj := range h.config.Projects {
		if proj.Name == name {
			h.config.Projects = append(h.config.Projects[:i], h.config.Projects[i+1:]...)
			if err := SaveConfig("config.json", h.config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to save configuration: %v", err),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Project deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Project not found",
	})
}

// Deploy initiates deployment for a project
func (h *ProjectHandler) Deploy(c *gin.Context) {
	name := c.Param("name")

	// Find the project
	var project *Project
	for _, proj := range h.config.Projects {
		if proj.Name == name {
			project = &proj
			break
		}
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Project not found",
		})
		return
	}

	// Validate that all deploy servers exist
	for _, serverName := range project.DeployServers {
		found := false
		for _, sshConfig := range h.config.SSHConfigs {
			if sshConfig.Name == serverName {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Deploy server '%s' not found in SSH configurations", serverName),
			})
			return
		}
	}

	// Generate a deployment ID
	deploymentID := fmt.Sprintf("%s-%d", project.Name, time.Now().Unix())

	// TODO: Actual deployment logic will be implemented in the deployment engine
	// For now, just return success with deployment ID

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"deploymentId": deploymentID,
		},
		"message": "Deployment initiated successfully",
	})
}
