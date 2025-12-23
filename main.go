package main

import (
	"embed"
	"flag"
	"log"

	"github.com/diiyw/ed/api"
	"github.com/diiyw/ed/api/handlers"
)

//go:embed frontend/dist
var embeddedFiles embed.FS

func main() {
	// Command-line flags
	apiMode := flag.Bool("api", false, "Run in API mode (web server)")
	port := flag.String("port", "8080", "API server port")
	flag.Parse()

	// Load or create config
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if *apiMode {
		// Run API server
		log.Printf("Starting API server on port %s...", *port)

		// Convert main.Config to handlers.Config
		handlerConfig := &handlers.Config{
			SSHConfigs: convertSSHConfigs(config.SSHConfigs),
			Projects:   convertProjects(config.Projects),
		}

		router := api.SetupRouter(handlerConfig, &embeddedFiles)
		if err := router.Run(":" + *port); err != nil {
			log.Fatal("Failed to start API server:", err)
		}
	} else {
		// TUI mode not implemented yet - use --api flag to run in API mode
		log.Fatal("TUI mode not available. Please use --api flag to run in API mode.")
	}
}

// Helper functions to convert between main and handlers types
func convertSSHConfigs(configs []SSHConfig) []handlers.SSHConfig {
	result := make([]handlers.SSHConfig, len(configs))
	for i, cfg := range configs {
		result[i] = handlers.SSHConfig{
			Name:     cfg.Name,
			Host:     cfg.Host,
			Port:     cfg.Port,
			User:     cfg.User,
			AuthType: cfg.AuthType,
			Password: cfg.Password,
			KeyFile:  cfg.KeyFile,
			KeyPass:  cfg.KeyPass,
		}
	}
	return result
}

func convertProjects(projects []Project) []handlers.Project {
	result := make([]handlers.Project, len(projects))
	for i, proj := range projects {
		result[i] = handlers.Project{
			Name:              proj.Name,
			BuildInstructions: proj.BuildInstructions,
			DeployScript:      proj.DeployScript,
			DeployServers:     proj.DeployServers,
			CreatedAt:         proj.CreatedAt,
			UpdatedAt:         proj.UpdatedAt,
		}
	}
	return result
}
