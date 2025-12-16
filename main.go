package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load or create config
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Start the TUI
	p := tea.NewProgram(NewMainMenu(config))
	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to run program:", err)
		os.Exit(1)
	}
}
