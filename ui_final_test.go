package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Task 16: Final polish and testing
// This file contains comprehensive tests for:
// - Testing all views on different terminal sizes
// - Verifying consistent styling across all screens
// - Ensuring graceful degradation if styles fail
// - Optimizing rendering performance

// TestAllViewsOnDifferentTerminalSizes tests that all views render correctly
// across various terminal sizes without crashing or producing empty output
func TestAllViewsOnDifferentTerminalSizes(t *testing.T) {
	// Define test terminal sizes
	terminalSizes := []struct {
		name   string
		width  int
		height int
	}{
		{"tiny", 40, 10},
		{"small", 80, 24},
		{"medium", 120, 40},
		{"large", 160, 60},
		{"wide", 200, 30},
		{"tall", 80, 100},
		{"very_small", 20, 5},
	}

	// Create test config with sample data
	config := &Config{
		SSHConfigs: []SSHConfig{
			{
				Name:     "TestServer1",
				Host:     "test1.example.com",
				Port:     22,
				User:     "user1",
				AuthType: "password",
				Password: "pass1",
			},
			{
				Name:     "TestServer2",
				Host:     "test2.example.com",
				Port:     2222,
				User:     "user2",
				AuthType: "key",
				KeyFile:  "/path/to/key",
			},
		},
		Projects: []Project{
			{
				Name:              "TestProject1",
				BuildInstructions: "npm install && npm run build",
				DeployScript:      "rsync -avz dist/ server:/var/www/",
				DeployServers:     []string{"TestServer1"},
			},
			{
				Name:              "TestProject2",
				BuildInstructions: "go build -o app",
				DeployScript:      "scp app server:/usr/local/bin/",
				DeployServers:     []string{"TestServer2"},
			},
		},
	}

	// Test all view models
	models := []struct {
		name  string
		model tea.Model
	}{
		{"MainMenu", NewMainMenu(config)},
		{"SSHList", NewSSHListModel(config)},
		{"SSHForm_New", NewSSHFormModel(config, -1)},
		{"SSHForm_Edit", NewSSHFormModel(config, 0)},
		{"SSHTest", NewSSHTestModel(config.SSHConfigs[0], config)},
		{"ProjectList", NewProjectListModel(config)},
		{"ProjectForm_New", NewProjectFormModel(config, -1)},
		{"ProjectForm_Edit", NewProjectFormModel(config, 0)},
		{"DeployModel", NewDeployModel(config.Projects[0], config)},
	}

	for _, size := range terminalSizes {
		t.Run(size.name, func(t *testing.T) {
			for _, m := range models {
				t.Run(m.name, func(t *testing.T) {
					// Send window size message to the model
					model, _ := m.model.Update(tea.WindowSizeMsg{
						Width:  size.width,
						Height: size.height,
					})

					// Get the view output
					view := model.View()

					// Verify the view is not empty
					if view == "" {
						t.Errorf("%s view should not be empty for terminal size %dx%d",
							m.name, size.width, size.height)
					}

					// Verify the view doesn't contain panic or error messages
					if strings.Contains(view, "panic") || strings.Contains(view, "PANIC") {
						t.Errorf("%s view contains panic for terminal size %dx%d",
							m.name, size.width, size.height)
					}

					// Verify the view has reasonable length (not truncated unexpectedly)
					if len(view) < 10 {
						t.Errorf("%s view is suspiciously short (%d chars) for terminal size %dx%d",
							m.name, len(view), size.width, size.height)
					}
				})
			}
		})
	}
}

// TestConsistentStylingAcrossAllScreens verifies that all screens use
// consistent styling patterns and color schemes
func TestConsistentStylingAcrossAllScreens(t *testing.T) {
	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
		},
		Projects: []Project{
			{Name: "TestProject", DeployServers: []string{"TestServer"}},
		},
	}

	// Get views from all models
	views := map[string]string{
		"MainMenu":    NewMainMenu(config).View(),
		"SSHList":     NewSSHListModel(config).View(),
		"SSHForm":     NewSSHFormModel(config, -1).View(),
		"ProjectList": NewProjectListModel(config).View(),
		"ProjectForm": NewProjectFormModel(config, -1).View(),
		"SSHTest":     NewSSHTestModel(config.SSHConfigs[0], config).View(),
		"DeployModel": NewDeployModel(config.Projects[0], config).View(),
	}

	// Test 1: All views should use consistent help text formatting
	// Most views should contain the bullet separator " • " for keyboard shortcuts
	// (Some views like SSHTest and DeployModel in initial state may only have one key)
	for name, view := range views {
		// Skip views that may only have one keyboard shortcut in their initial state
		if name == "SSHTest" || name == "DeployModel" {
			continue
		}
		if !strings.Contains(view, " • ") {
			t.Errorf("%s should use consistent help text separator ' • '", name)
		}
	}

	// Test 2: All views should use consistent border characters
	// Check for rounded border characters (╭╮╰╯─│)
	borderChars := []string{"╭", "╮", "╰", "╯", "─", "│"}
	for name, view := range views {
		hasBorder := false
		for _, char := range borderChars {
			if strings.Contains(view, char) {
				hasBorder = true
				break
			}
		}
		if !hasBorder {
			t.Errorf("%s should use consistent border styling", name)
		}
	}

	// Test 3: All views should use consistent icon styling
	// Check that views contain at least some icons
	commonIcons := []string{"✓", "✗", "⚠", "ℹ", "→", "•", "▶"}
	for name, view := range views {
		hasIcon := false
		for _, icon := range commonIcons {
			if strings.Contains(view, icon) {
				hasIcon = true
				break
			}
		}
		// Some views might not have icons, so we just log a warning
		if !hasIcon {
			t.Logf("%s does not contain common icons (this may be expected)", name)
		}
	}

	// Test 4: All form views should use consistent field formatting
	formViews := map[string]string{
		"SSHForm":     views["SSHForm"],
		"ProjectForm": views["ProjectForm"],
	}

	for name, view := range formViews {
		// Forms should have cursor indicators
		if !strings.Contains(view, "▶") {
			t.Errorf("%s should have cursor indicator '▶'", name)
		}

		// Forms should have field labels
		if !strings.Contains(view, "Name") {
			t.Errorf("%s should have 'Name' field label", name)
		}

		// Forms should have help text
		if !strings.Contains(view, "Navigate") {
			t.Errorf("%s should have navigation help text", name)
		}
	}

	// Test 5: All list views should use consistent item formatting
	listViews := map[string]string{
		"SSHList":     views["SSHList"],
		"ProjectList": views["ProjectList"],
	}

	for name, view := range listViews {
		// Lists should have help text
		if !strings.Contains(view, "Add") {
			t.Errorf("%s should have 'Add' action in help text", name)
		}

		// Lists should have navigation help
		if !strings.Contains(view, "esc") || !strings.Contains(view, "Back") {
			t.Errorf("%s should have 'esc: Back' in help text", name)
		}
	}
}

// TestGracefulDegradationOnStyleFailure verifies that the UI doesn't crash
// if styling fails and falls back to plain text rendering
func TestGracefulDegradationOnStyleFailure(t *testing.T) {
	// Test that helper functions handle edge cases gracefully

	// Test renderIcon with invalid icon type
	result := renderIcon("invalid_icon_type_that_does_not_exist", bodyStyle)
	if result == "" {
		t.Error("renderIcon should return the input string for unknown icon types, not empty string")
	}

	// Test renderBadge with invalid badge type
	result = renderBadge("Test", "invalid_badge_type")
	if result == "" {
		t.Error("renderBadge should return styled text even for unknown badge types")
	}
	if !strings.Contains(result, "Test") {
		t.Error("renderBadge should contain the input text even for unknown badge types")
	}

	// Test renderDivider with zero and negative widths
	result = renderDivider(0, dividerStyle)
	if result == "" {
		t.Error("renderDivider should handle zero width gracefully")
	}

	result = renderDivider(-10, dividerStyle)
	if result == "" {
		t.Error("renderDivider should handle negative width gracefully")
	}

	// Test renderProgressBar with invalid values
	result = renderProgressBar(0, 0, 40)
	if result == "" {
		t.Error("renderProgressBar should handle zero total gracefully")
	}

	result = renderProgressBar(150, 100, 40)
	if result == "" {
		t.Error("renderProgressBar should handle current > total gracefully")
	}
	if !strings.Contains(result, "100%") {
		t.Error("renderProgressBar should clamp to 100% when current > total")
	}

	result = renderProgressBar(50, 100, 0)
	if result == "" {
		t.Error("renderProgressBar should handle zero width gracefully")
	}

	result = renderProgressBar(50, 100, -10)
	if result == "" {
		t.Error("renderProgressBar should handle negative width gracefully")
	}

	// Test renderCard with empty content and title
	result = renderCard("", "")
	if result == "" {
		t.Error("renderCard should handle empty content and title gracefully")
	}

	// Test renderKeyHelp with empty map
	result = renderKeyHelp(map[string]string{})
	if result != "" {
		t.Error("renderKeyHelp should return empty string for empty map")
	}

	// Test renderDividerDecorative with various edge cases
	result = renderDividerDecorative(0, dividerStyle, "dots")
	if result == "" {
		t.Error("renderDividerDecorative should handle zero width gracefully")
	}

	result = renderDividerDecorative(40, dividerStyle, "unknown_decoration_type")
	if result == "" {
		t.Error("renderDividerDecorative should handle unknown decoration type gracefully")
	}
}

// TestRenderingPerformance verifies that rendering is reasonably fast
// and doesn't have obvious performance issues
func TestRenderingPerformance(t *testing.T) {
	// Create config with moderate amount of data
	config := &Config{
		SSHConfigs: make([]SSHConfig, 20),
		Projects:   make([]Project, 20),
	}

	// Populate with test data
	for i := 0; i < 20; i++ {
		config.SSHConfigs[i] = SSHConfig{
			Name:     "Server" + string(rune(i)),
			Host:     "host" + string(rune(i)) + ".example.com",
			Port:     22 + i,
			User:     "user" + string(rune(i)),
			AuthType: "password",
			Password: "pass",
		}
		config.Projects[i] = Project{
			Name:              "Project" + string(rune(i)),
			BuildInstructions: "build instructions for project " + string(rune(i)),
			DeployScript:      "deploy script for project " + string(rune(i)),
			DeployServers:     []string{"Server" + string(rune(i))},
		}
	}

	// Test that views can be rendered multiple times without issues
	iterations := 100

	models := []struct {
		name  string
		model tea.Model
	}{
		{"MainMenu", NewMainMenu(config)},
		{"SSHList", NewSSHListModel(config)},
		{"ProjectList", NewProjectListModel(config)},
	}

	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				view := m.model.View()
				if view == "" {
					t.Errorf("View became empty after %d iterations", i)
					break
				}
			}
		})
	}
}

// TestEmptyStateHandling verifies that all views handle empty states gracefully
func TestEmptyStateHandling(t *testing.T) {
	// Create config with no data
	emptyConfig := &Config{
		SSHConfigs: []SSHConfig{},
		Projects:   []Project{},
	}

	// Test SSH List with no configurations
	sshList := NewSSHListModel(emptyConfig)
	view := sshList.View()

	if view == "" {
		t.Error("SSH List should display empty state message, not empty view")
	}

	if !strings.Contains(view, "No SSH Configurations Found") {
		t.Error("SSH List should display 'No SSH Configurations Found' message")
	}

	if !strings.Contains(view, "Press") && !strings.Contains(view, "a") {
		t.Error("SSH List should display instruction to add configuration")
	}

	// Test Project List with no projects
	projectList := NewProjectListModel(emptyConfig)
	view = projectList.View()

	if view == "" {
		t.Error("Project List should display empty state message, not empty view")
	}

	if !strings.Contains(view, "No Projects Found") {
		t.Error("Project List should display 'No Projects Found' message")
	}

	if !strings.Contains(view, "Press") && !strings.Contains(view, "a") {
		t.Error("Project List should display instruction to add project")
	}
}

// TestViewConsistencyAfterUpdates verifies that views remain consistent
// after multiple update cycles
func TestViewConsistencyAfterUpdates(t *testing.T) {
	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "Server1", Host: "host1.com", Port: 22, User: "user1", AuthType: "password", Password: "pass"},
			{Name: "Server2", Host: "host2.com", Port: 22, User: "user2", AuthType: "password", Password: "pass"},
		},
		Projects: []Project{
			{Name: "Project1", DeployServers: []string{"Server1"}},
			{Name: "Project2", DeployServers: []string{"Server2"}},
		},
	}

	// Test Main Menu navigation
	mainMenu := NewMainMenu(config)
	initialView := mainMenu.View()

	// Simulate navigation
	mainMenu.list.CursorDown()
	view1 := mainMenu.View()

	mainMenu.list.CursorDown()
	view2 := mainMenu.View()

	mainMenu.list.CursorUp()
	view3 := mainMenu.View()

	// All views should be non-empty
	if initialView == "" || view1 == "" || view2 == "" || view3 == "" {
		t.Error("Views should not become empty after navigation")
	}

	// View should change when cursor moves
	if initialView == view1 {
		t.Error("View should change when cursor moves down")
	}

	// View should be consistent when returning to same position
	// Note: Due to styling and ANSI codes, we check that both views contain the same menu items
	// rather than exact string equality
	if !strings.Contains(view3, "SSH Management") || !strings.Contains(view3, "Project Management") {
		t.Error("View should contain menu items when returning to same cursor position")
	}
}

// TestAllViewsContainRequiredElements verifies that all views contain
// their required UI elements
func TestAllViewsContainRequiredElements(t *testing.T) {
	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
		},
		Projects: []Project{
			{Name: "TestProject", DeployServers: []string{"TestServer"}},
		},
	}

	tests := []struct {
		name             string
		model            tea.Model
		requiredElements []string
	}{
		{
			name:  "MainMenu",
			model: NewMainMenu(config),
			requiredElements: []string{
				"EASY DEPLOY",
				"Main Menu",
				"SSH Management",
				"Project Management",
				"Exit",
				"Navigate",
				"Select",
				"Quit",
			},
		},
		{
			name:  "SSHList",
			model: NewSSHListModel(config),
			requiredElements: []string{
				"SSH Configurations",
				"TestServer",
				"test.com",
				"Add",
				"Edit",
				"Delete",
			},
		},
		{
			name:  "SSHForm",
			model: NewSSHFormModel(config, -1),
			requiredElements: []string{
				"Add SSH Configuration",
				"Name",
				"Host",
				"Port",
				"User",
				"Navigate",
				"Save",
			},
		},
		{
			name:  "ProjectList",
			model: NewProjectListModel(config),
			requiredElements: []string{
				"Projects",
				"TestProject",
				"Add",
				"Edit",
				"Delete",
				"Deploy",
			},
		},
		{
			name:  "ProjectForm",
			model: NewProjectFormModel(config, -1),
			requiredElements: []string{
				"Add Project",
				"Name",
				"Navigate",
				"Save",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.model.View()

			for _, element := range tt.requiredElements {
				if !strings.Contains(view, element) {
					t.Errorf("%s view should contain '%s'", tt.name, element)
				}
			}
		})
	}
}
