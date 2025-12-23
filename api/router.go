package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/diiyw/ed/api/handlers"
	"github.com/diiyw/ed/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all API routes
func SetupRouter(config *handlers.Config, embeddedFS *embed.FS) *gin.Engine {
	// Create router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.SetupCORS())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())

	// Create handlers
	sshHandler := handlers.NewSSHHandler(config)
	projectHandler := handlers.NewProjectHandler(config)

	// API routes
	api := router.Group("/api")
	{
		// SSH configuration routes
		ssh := api.Group("/ssh")
		{
			ssh.GET("", sshHandler.GetAll)
			ssh.GET("/:name", sshHandler.GetByName)
			ssh.POST("", sshHandler.Create)
			ssh.PUT("/:name", sshHandler.Update)
			ssh.DELETE("/:name", sshHandler.Delete)
			ssh.POST("/:name/test", sshHandler.Test)
		}

		// Project routes
		projects := api.Group("/projects")
		{
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:name", projectHandler.GetByName)
			projects.POST("", projectHandler.Create)
			projects.PUT("/:name", projectHandler.Update)
			projects.DELETE("/:name", projectHandler.Delete)
			projects.POST("/:name/deploy", projectHandler.Deploy)
		}
	}

	// WebSocket routes
	// TODO: Add WebSocket handler for deployment logs
	// router.GET("/ws/deploy/:deploymentId", websocketHandler.HandleDeployment)

	// Serve embedded frontend files if provided
	if embeddedFS != nil {
		setupStaticRoutes(router, embeddedFS)
	}

	return router
}

// setupStaticRoutes configures routes to serve the embedded frontend
func setupStaticRoutes(router *gin.Engine, embeddedFS *embed.FS) {
	// Try to get the embedded filesystem starting from frontend/dist
	distFS, err := fs.Sub(*embeddedFS, "frontend/dist")
	if err != nil {
		// Try testdata/dist for tests
		distFS, err = fs.Sub(*embeddedFS, "testdata/dist")
		if err != nil {
			// If we can't get either subdirectory, just return
			return
		}
	}

	// Serve static files
	router.NoRoute(func(c *gin.Context) {
		// Skip if this is an API or WebSocket route
		path := c.Request.URL.Path
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		if len(path) >= 3 && path[:3] == "/ws" {
			c.JSON(http.StatusNotFound, gin.H{"error": "WebSocket endpoint not found"})
			return
		}

		// Try to serve the requested file
		filePath := path
		if filePath == "/" {
			filePath = "index.html"
		} else if len(filePath) > 0 && filePath[0] == '/' {
			// Remove leading slash for fs.FS
			filePath = filePath[1:]
		}

		// Try to read the file
		content, err := fs.ReadFile(distFS, filePath)
		if err != nil {
			// File not found, serve index.html for SPA routing
			filePath = "index.html"
			content, err = fs.ReadFile(distFS, filePath)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
				return
			}
		}

		// Detect content type based on file extension
		contentType := "text/html; charset=utf-8"
		if len(filePath) > 3 && filePath[len(filePath)-3:] == ".js" {
			contentType = "application/javascript; charset=utf-8"
		} else if len(filePath) > 4 && filePath[len(filePath)-4:] == ".css" {
			contentType = "text/css; charset=utf-8"
		} else if len(filePath) > 4 && filePath[len(filePath)-4:] == ".svg" {
			contentType = "image/svg+xml"
		}

		// Set headers and write content
		c.Data(http.StatusOK, contentType, content)
	})
}
