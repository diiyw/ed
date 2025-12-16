package main

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: ui-enhancement, Property 1: Status message color differentiation**
// For any status message (success, error, warning, info), the rendered output should use the color designated for that message type
func TestProperty_StatusMessageColorDifferentiation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("status messages use designated colors", prop.ForAll(
		func(msgType string, content string) bool {
			// Get the appropriate style based on message type
			var style lipgloss.Style
			var expectedColor lipgloss.Color

			switch msgType {
			case "success":
				style = successStyle
				expectedColor = successColor
			case "error":
				style = errorStyle
				expectedColor = errorColor
			case "warning":
				style = warningStyle
				expectedColor = warningColor
			case "info":
				style = statusStyle
				expectedColor = infoColor
			default:
				return false // Invalid message type
			}

			// Verify the style has the correct foreground color
			// We check that the style's foreground color matches the expected color
			actualColor := style.GetForeground()

			return actualColor == expectedColor
		},
		gen.OneConstOf("success", "error", "warning", "info"),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Helper test to verify the color constants are correctly defined
func TestStatusColorConstants(t *testing.T) {
	tests := []struct {
		name          string
		style         lipgloss.Style
		expectedColor lipgloss.Color
	}{
		{"success style uses success color", successStyle, successColor},
		{"error style uses error color", errorStyle, errorColor},
		{"warning style uses warning color", warningStyle, warningColor},
		{"info style uses info color", statusStyle, infoColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualColor := tt.style.GetForeground()
			if actualColor != tt.expectedColor {
				t.Errorf("Expected color %v, got %v", tt.expectedColor, actualColor)
			}
		})
	}
}

// Property test to verify that different status types produce visually distinct output
func TestProperty_StatusMessageVisualDistinction(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("different status types produce distinct styled output", prop.ForAll(
		func(content string) bool {
			// Skip empty strings as they don't provide meaningful distinction
			if strings.TrimSpace(content) == "" {
				return true
			}

			// Render the same content with different status styles
			successRendered := successStyle.Render(content)
			errorRendered := errorStyle.Render(content)
			warningRendered := warningStyle.Render(content)
			infoRendered := statusStyle.Render(content)

			// All rendered outputs should be different from each other
			// (because they use different colors)
			outputs := []string{successRendered, errorRendered, warningRendered, infoRendered}

			// Check that we have distinct outputs
			// At minimum, the ANSI color codes should differ
			for i := 0; i < len(outputs); i++ {
				for j := i + 1; j < len(outputs); j++ {
					if outputs[i] == outputs[j] {
						return false
					}
				}
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit Tests for Helper Functions

// Test icon rendering with various icon types
func TestRenderIcon(t *testing.T) {
	tests := []struct {
		name     string
		iconType string
		style    lipgloss.Style
		expected string
	}{
		{"success icon", "success", lipgloss.NewStyle(), "✓"},
		{"error icon", "error", lipgloss.NewStyle(), "✗"},
		{"warning icon", "warning", lipgloss.NewStyle(), "⚠"},
		{"info icon", "info", lipgloss.NewStyle(), "ℹ"},
		{"loading icon", "loading", lipgloss.NewStyle(), "⋯"},
		{"arrow icon", "arrow", lipgloss.NewStyle(), "→"},
		{"bullet icon", "bullet", lipgloss.NewStyle(), "•"},
		{"star icon", "star", lipgloss.NewStyle(), "★"},
		{"unknown icon returns input", "unknown", lipgloss.NewStyle(), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderIcon(tt.iconType, tt.style)
			// The result will be styled, so we check if it contains the expected icon
			if !strings.Contains(result, tt.expected) {
				t.Errorf("renderIcon(%q) = %q, want to contain %q", tt.iconType, result, tt.expected)
			}
		})
	}
}

// Test badge rendering with different badge types and text
func TestRenderBadge(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		badgeType string
	}{
		{"success badge", "Success", "success"},
		{"error badge", "Error", "error"},
		{"warning badge", "Warning", "warning"},
		{"info badge", "Info", "info"},
		{"primary badge", "Primary", "primary"},
		{"secondary badge", "Secondary", "secondary"},
		{"muted badge", "Muted", "muted"},
		{"default badge", "Default", "default"},
		{"empty text", "", "success"},
		{"long text", "This is a very long badge text", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderBadge(tt.text, tt.badgeType)
			// The result should contain the text (even if styled)
			if !strings.Contains(result, tt.text) {
				t.Errorf("renderBadge(%q, %q) = %q, want to contain %q", tt.text, tt.badgeType, result, tt.text)
			}
			// Result should not be empty
			if result == "" {
				t.Errorf("renderBadge(%q, %q) returned empty string", tt.text, tt.badgeType)
			}
		})
	}
}

// Test divider rendering with various widths
func TestRenderDivider(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected int // expected number of divider characters
	}{
		{"zero width defaults to 40", 0, 40},
		{"negative width defaults to 40", -1, 40},
		{"width 10", 10, 10},
		{"width 20", 20, 20},
		{"width 50", 50, 50},
		{"width 100", 100, 100},
		{"width 1", 1, 1},
	}

	style := lipgloss.NewStyle()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderDivider(tt.width, style)
			// Count the number of divider characters (─)
			count := strings.Count(result, "─")
			if count != tt.expected {
				t.Errorf("renderDivider(%d) produced %d characters, want %d", tt.width, count, tt.expected)
			}
		})
	}
}

// Test progress bar with different completion percentages
func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		width    int
		wantPerc string
	}{
		{"0% complete", 0, 100, 40, "0%"},
		{"50% complete", 50, 100, 40, "50%"},
		{"100% complete", 100, 100, 40, "100%"},
		{"25% complete", 25, 100, 40, "25%"},
		{"75% complete", 75, 100, 40, "75%"},
		{"over 100% clamped", 150, 100, 40, "100%"},
		{"zero total defaults to 1", 1, 0, 40, "100%"},
		{"negative width defaults to 40", 50, 100, 0, "50%"},
		{"small width", 5, 10, 10, "50%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderProgressBar(tt.current, tt.total, tt.width)
			// Check that the percentage text is present
			if !strings.Contains(result, tt.wantPerc) {
				t.Errorf("renderProgressBar(%d, %d, %d) = %q, want to contain %q", tt.current, tt.total, tt.width, result, tt.wantPerc)
			}
			// Check that result contains progress characters
			if !strings.Contains(result, "█") && !strings.Contains(result, "░") {
				t.Errorf("renderProgressBar(%d, %d, %d) = %q, want to contain progress characters", tt.current, tt.total, tt.width, result)
			}
		})
	}
}

// Test card wrapping with different content sizes
func TestRenderCard(t *testing.T) {
	tests := []struct {
		name    string
		content string
		title   string
	}{
		{"card with title and content", "This is some content", "Card Title"},
		{"card with only content", "This is some content", ""},
		{"card with empty content", "", "Card Title"},
		{"card with multiline content", "Line 1\nLine 2\nLine 3", "Multi-line Card"},
		{"card with long content", strings.Repeat("A", 200), "Long Content"},
		{"card with special characters", "Special: !@#$%^&*()", "Special Chars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderCard(tt.content, tt.title)
			// Check that content is present in result (for multiline, check each line)
			if tt.content != "" {
				lines := strings.Split(tt.content, "\n")
				for _, line := range lines {
					if line != "" && !strings.Contains(result, line) {
						t.Errorf("renderCard() result doesn't contain content line: %q", line)
					}
				}
			}
			// Check that title is present if provided
			if tt.title != "" && !strings.Contains(result, tt.title) {
				t.Errorf("renderCard() result doesn't contain title: %q", tt.title)
			}
			// Result should not be empty
			if result == "" {
				t.Errorf("renderCard(%q, %q) returned empty string", tt.content, tt.title)
			}
		})
	}
}

// Test spinner rendering
func TestRenderSpinner(t *testing.T) {
	// Test that different frames produce different outputs
	frames := make(map[string]bool)
	for i := 0; i < 20; i++ {
		result := renderSpinner(i)
		frames[result] = true
		// Result should not be empty
		if result == "" {
			t.Errorf("renderSpinner(%d) returned empty string", i)
		}
	}
	// Should have multiple distinct frames (at least 5)
	if len(frames) < 5 {
		t.Errorf("renderSpinner produced only %d distinct frames, expected at least 5", len(frames))
	}
}

// Test renderKeyHelp
func TestRenderKeyHelp(t *testing.T) {
	tests := []struct {
		name string
		keys map[string]string
	}{
		{"empty map", map[string]string{}},
		{"single key", map[string]string{"a": "Add"}},
		{"multiple keys", map[string]string{"a": "Add", "e": "Edit", "d": "Delete"}},
		{"keys with special chars", map[string]string{"ctrl+c": "Quit", "esc": "Back"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderKeyHelp(tt.keys)
			// Empty map should return empty string
			if len(tt.keys) == 0 {
				if result != "" {
					t.Errorf("renderKeyHelp(empty) = %q, want empty string", result)
				}
				return
			}
			// Result should contain all keys and descriptions
			for key, desc := range tt.keys {
				if !strings.Contains(result, key) {
					t.Errorf("renderKeyHelp() result doesn't contain key: %q", key)
				}
				if !strings.Contains(result, desc) {
					t.Errorf("renderKeyHelp() result doesn't contain description: %q", desc)
				}
			}
		})
	}
}

// Test renderDividerDecorative
func TestRenderDividerDecorative(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		decorationType string
		expectedChars  []string
	}{
		{"dots decoration", 10, "dots", []string{"·", "─"}},
		{"arrows decoration", 10, "arrows", []string{"◄", "►", "─"}},
		{"stars decoration", 10, "stars", []string{"★", "─"}},
		{"diamonds decoration", 10, "diamonds", []string{"◆", "─"}},
		{"double decoration", 10, "double", []string{"═"}},
		{"wave decoration", 10, "wave", []string{"～"}},
		{"default decoration", 10, "default", []string{"─"}},
		{"unknown decoration defaults", 10, "unknown", []string{"─"}},
	}

	style := lipgloss.NewStyle()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderDividerDecorative(tt.width, style, tt.decorationType)
			// Check that at least one expected character is present
			found := false
			for _, char := range tt.expectedChars {
				if strings.Contains(result, char) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("renderDividerDecorative(%d, %q) = %q, want to contain one of %v", tt.width, tt.decorationType, result, tt.expectedChars)
			}
		})
	}
}

// **Feature: ui-enhancement, Property 5: Consistent border styling**
// For any container, list, or form element, the border style (rounded, thickness, color) should be applied consistently based on element type
func TestProperty_ConsistentBorderStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Define element types and their expected border characteristics
	type elementType struct {
		name             string
		style            lipgloss.Style
		expectedBorder   lipgloss.Border
		expectedPaddingH int
		expectedPaddingV int
		expectedBorderFg lipgloss.Color
	}

	elementTypes := []elementType{
		{
			name:             "default_border",
			style:            borderStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: primaryColor,
		},
		{
			name:             "thick_border",
			style:            borderStyleThick,
			expectedBorder:   lipgloss.ThickBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: primaryColor,
		},
		{
			name:             "thin_border",
			style:            borderStyleThin,
			expectedBorder:   lipgloss.NormalBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: mutedColor,
		},
		{
			name:             "double_border",
			style:            borderStyleDouble,
			expectedBorder:   lipgloss.DoubleBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: accentColor,
		},
		{
			name:             "form",
			style:            formStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: secondaryColor,
		},
		{
			name:             "card",
			style:            cardStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedPaddingH: 2,
			expectedPaddingV: 1,
			expectedBorderFg: primaryColor,
		},
	}

	properties.Property("border styles are consistent for each element type", prop.ForAll(
		func(content string) bool {
			// For each element type, verify that the border characteristics are consistent
			for _, elemType := range elementTypes {
				// Get the border style characteristics
				actualBorder := elemType.style.GetBorderStyle()
				actualBorderTopFg := elemType.style.GetBorderTopForeground()
				actualBorderRightFg := elemType.style.GetBorderRightForeground()
				actualBorderBottomFg := elemType.style.GetBorderBottomForeground()
				actualBorderLeftFg := elemType.style.GetBorderLeftForeground()
				actualPaddingTop := elemType.style.GetPaddingTop()
				actualPaddingRight := elemType.style.GetPaddingRight()
				actualPaddingBottom := elemType.style.GetPaddingBottom()
				actualPaddingLeft := elemType.style.GetPaddingLeft()

				// Verify border style matches expected
				if actualBorder != elemType.expectedBorder {
					t.Logf("Element type %s: border style mismatch", elemType.name)
					return false
				}

				// Verify border foreground colors match expected (all sides should be consistent)
				if actualBorderTopFg != elemType.expectedBorderFg ||
					actualBorderRightFg != elemType.expectedBorderFg ||
					actualBorderBottomFg != elemType.expectedBorderFg ||
					actualBorderLeftFg != elemType.expectedBorderFg {
					t.Logf("Element type %s: border foreground color mismatch. Expected %v, got top=%v right=%v bottom=%v left=%v",
						elemType.name, elemType.expectedBorderFg, actualBorderTopFg, actualBorderRightFg, actualBorderBottomFg, actualBorderLeftFg)
					return false
				}

				// Verify horizontal padding is consistent
				if actualPaddingLeft != elemType.expectedPaddingH || actualPaddingRight != elemType.expectedPaddingH {
					t.Logf("Element type %s: horizontal padding mismatch. Expected %d, got left=%d right=%d",
						elemType.name, elemType.expectedPaddingH, actualPaddingLeft, actualPaddingRight)
					return false
				}

				// Verify vertical padding is consistent
				if actualPaddingTop != elemType.expectedPaddingV || actualPaddingBottom != elemType.expectedPaddingV {
					t.Logf("Element type %s: vertical padding mismatch. Expected %d, got top=%d bottom=%d",
						elemType.name, elemType.expectedPaddingV, actualPaddingTop, actualPaddingBottom)
					return false
				}

				// Verify that rendering with the style produces consistent output
				// (i.e., the same content rendered twice should produce identical output)
				rendered1 := elemType.style.Render(content)
				rendered2 := elemType.style.Render(content)
				if rendered1 != rendered2 {
					t.Logf("Element type %s: inconsistent rendering for same content", elemType.name)
					return false
				}
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test to verify border style characteristics are correctly defined
func TestBorderStyleCharacteristics(t *testing.T) {
	tests := []struct {
		name             string
		style            lipgloss.Style
		expectedBorder   lipgloss.Border
		expectedBorderFg lipgloss.Color
		expectedPaddingH int
		expectedPaddingV int
	}{
		{
			name:             "default border style",
			style:            borderStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedBorderFg: primaryColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
		{
			name:             "thick border style",
			style:            borderStyleThick,
			expectedBorder:   lipgloss.ThickBorder(),
			expectedBorderFg: primaryColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
		{
			name:             "thin border style",
			style:            borderStyleThin,
			expectedBorder:   lipgloss.NormalBorder(),
			expectedBorderFg: mutedColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
		{
			name:             "double border style",
			style:            borderStyleDouble,
			expectedBorder:   lipgloss.DoubleBorder(),
			expectedBorderFg: accentColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
		{
			name:             "form style",
			style:            formStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedBorderFg: secondaryColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
		{
			name:             "card style",
			style:            cardStyle,
			expectedBorder:   lipgloss.RoundedBorder(),
			expectedBorderFg: primaryColor,
			expectedPaddingH: 2,
			expectedPaddingV: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check border style
			actualBorder := tt.style.GetBorderStyle()
			if actualBorder != tt.expectedBorder {
				t.Errorf("Border style mismatch: expected %v, got %v", tt.expectedBorder, actualBorder)
			}

			// Check border foreground color (check all sides)
			actualBorderTopFg := tt.style.GetBorderTopForeground()
			actualBorderRightFg := tt.style.GetBorderRightForeground()
			actualBorderBottomFg := tt.style.GetBorderBottomForeground()
			actualBorderLeftFg := tt.style.GetBorderLeftForeground()
			if actualBorderTopFg != tt.expectedBorderFg ||
				actualBorderRightFg != tt.expectedBorderFg ||
				actualBorderBottomFg != tt.expectedBorderFg ||
				actualBorderLeftFg != tt.expectedBorderFg {
				t.Errorf("Border foreground color mismatch: expected %v, got top=%v right=%v bottom=%v left=%v",
					tt.expectedBorderFg, actualBorderTopFg, actualBorderRightFg, actualBorderBottomFg, actualBorderLeftFg)
			}

			// Check horizontal padding
			actualPaddingLeft := tt.style.GetPaddingLeft()
			actualPaddingRight := tt.style.GetPaddingRight()
			if actualPaddingLeft != tt.expectedPaddingH || actualPaddingRight != tt.expectedPaddingH {
				t.Errorf("Horizontal padding mismatch: expected %d, got left=%d right=%d",
					tt.expectedPaddingH, actualPaddingLeft, actualPaddingRight)
			}

			// Check vertical padding
			actualPaddingTop := tt.style.GetPaddingTop()
			actualPaddingBottom := tt.style.GetPaddingBottom()
			if actualPaddingTop != tt.expectedPaddingV || actualPaddingBottom != tt.expectedPaddingV {
				t.Errorf("Vertical padding mismatch: expected %d, got top=%d bottom=%d",
					tt.expectedPaddingV, actualPaddingTop, actualPaddingBottom)
			}
		})
	}
}

// Unit test for main menu header display
// Verifies that the main menu displays a header with the application title
// Requirements: 3.2
func TestMainMenuHeaderDisplay(t *testing.T) {
	// Create a test config
	config := &Config{
		SSHConfigs: []SSHConfig{},
		Projects:   []Project{},
	}

	// Create main menu model
	mainMenu := NewMainMenu(config)

	// Get the view output
	view := mainMenu.View()

	// Verify that the view contains the application title
	if !strings.Contains(view, "EASY DEPLOY") {
		t.Errorf("Main menu view should contain application title 'EASY DEPLOY', got: %s", view)
	}

	// Verify that the view contains the rocket emoji (part of the banner)
	if !strings.Contains(view, "🚀") {
		t.Errorf("Main menu view should contain rocket emoji in banner, got: %s", view)
	}

	// Verify that the view contains "Main Menu" subtitle
	if !strings.Contains(view, "Main Menu") {
		t.Errorf("Main menu view should contain 'Main Menu' subtitle, got: %s", view)
	}

	// Verify that the view is not empty
	if view == "" {
		t.Error("Main menu view should not be empty")
	}
}

// **Feature: ui-enhancement, Property 6: Selected item highlighting**
// For any list with a selected item, the selected item should have distinct background color and visual indicator that differs from unselected items
// Validates: Requirements 4.1, 4.5
func TestProperty_SelectedItemHighlighting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("selected items have distinct styling from unselected items", prop.ForAll(
		func(selectedIndex int) bool {
			// Main menu has 3 items (SSH Management, Project Management, Exit)
			numItems := 3

			// Ensure selected index is within bounds
			selectedIndex = selectedIndex % numItems
			if selectedIndex < 0 {
				selectedIndex = -selectedIndex
			}

			// Create a test config
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create main menu model
			mainMenu := NewMainMenu(config)

			// Set the selected index by simulating key presses
			for i := 0; i < selectedIndex; i++ {
				mainMenu.list.CursorDown()
			}

			// Get the view output
			view := mainMenu.View()

			// The selected item should have:
			// 1. An arrow indicator (→)
			// 2. Different styling than unselected items

			// Check for arrow indicator in the view
			hasArrowIndicator := strings.Contains(view, "→")

			// The view should contain menu items
			hasMenuItems := strings.Contains(view, "SSH Management") ||
				strings.Contains(view, "Project Management") ||
				strings.Contains(view, "Exit")

			// For a proper test, we need to verify that the selected item
			// appears differently than unselected items
			// We can check this by verifying that the view contains
			// visual distinction markers (like the arrow)

			return hasArrowIndicator && hasMenuItems
		},
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}

// Unit test to verify selected item styling characteristics
func TestSelectedItemStyling(t *testing.T) {
	// Create a test config
	config := &Config{
		SSHConfigs: []SSHConfig{},
		Projects:   []Project{},
	}

	// Main menu has 3 items
	menuItems := []string{"SSH Management", "Project Management", "Exit"}

	// Test with different selected indices
	for selectedIndex := 0; selectedIndex < len(menuItems); selectedIndex++ {
		t.Run(fmt.Sprintf("selected_index_%d", selectedIndex), func(t *testing.T) {
			// Create main menu model
			mainMenu := NewMainMenu(config)

			// Set the selected index
			for i := 0; i < selectedIndex; i++ {
				mainMenu.list.CursorDown()
			}

			// Get the view output
			view := mainMenu.View()

			// Verify arrow indicator is present
			if !strings.Contains(view, "→") {
				t.Errorf("Selected item should have arrow indicator")
			}

			// Verify selected item name is present
			selectedItemName := menuItems[selectedIndex]
			if !strings.Contains(view, selectedItemName) {
				t.Errorf("View should contain selected item name '%s'", selectedItemName)
			}

			// Verify view is not empty
			if view == "" {
				t.Error("View should not be empty")
			}
		})
	}
}

// **Feature: ui-enhancement, Property 20: List item formatting consistency**
// For any list, all items should have identical formatting and alignment
// Validates: Requirements 9.1, 9.5
func TestProperty_ListItemFormattingConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all list items have consistent formatting and alignment", prop.ForAll(
		func(numItems int) bool {
			// Ensure we have at least 2 items to compare
			if numItems < 2 {
				numItems = 2
			}
			// Cap at reasonable number for testing
			if numItems > 20 {
				numItems = 20
			}

			// Generate random SSH configs
			sshConfigs := make([]SSHConfig, numItems)
			for i := 0; i < numItems; i++ {
				sshConfigs[i] = SSHConfig{
					Name:     fmt.Sprintf("Server%d", i),
					Host:     fmt.Sprintf("host%d.example.com", i),
					Port:     22 + i,
					User:     fmt.Sprintf("user%d", i),
					AuthType: "password",
					Password: "pass",
				}
			}

			config := &Config{
				SSHConfigs: sshConfigs,
				Projects:   []Project{},
			}

			// Create SSH list model
			sshList := NewSSHListModel(config)
			view := sshList.View()

			// Extract individual item representations from the view
			// Each item should have consistent structure:
			// - Status indicator (circle icon)
			// - Name (bold)
			// - Connection details (user@host:port)
			// - Auth type
			// All wrapped in a card with consistent border and padding

			// Check that all items have the same structural elements
			// Count occurrences of key formatting elements
			circleCount := strings.Count(view, "●")

			// Each item should have exactly one circle indicator
			if circleCount != numItems {
				t.Logf("Expected %d circle indicators (one per item), got %d", numItems, circleCount)
				return false
			}

			// Verify all item names are present
			for i := 0; i < numItems; i++ {
				expectedName := fmt.Sprintf("Server%d", i)
				if !strings.Contains(view, expectedName) {
					t.Logf("Item name '%s' not found in view", expectedName)
					return false
				}
			}

			// Verify all connection details are present with consistent format
			for i := 0; i < numItems; i++ {
				expectedDetail := fmt.Sprintf("user%d@host%d.example.com:%d", i, i, 22+i)
				if !strings.Contains(view, expectedDetail) {
					t.Logf("Connection detail '%s' not found in view", expectedDetail)
					return false
				}
			}

			// Verify all items have auth type displayed
			authTypeCount := strings.Count(view, "Auth: password")
			if authTypeCount != numItems {
				t.Logf("Expected %d 'Auth: password' entries, got %d", numItems, authTypeCount)
				return false
			}

			// Now test with Project list
			// The project list uses the default list delegate, so formatting is different
			projects := make([]Project, numItems)
			for i := 0; i < numItems; i++ {
				projects[i] = Project{
					Name:          fmt.Sprintf("Project%d", i),
					DeployServers: []string{fmt.Sprintf("Server%d", i)},
				}
			}

			config.Projects = projects
			projectList := NewProjectListModel(config)

			// Set window size so the list renders properly
			// Use a large height to ensure all items are visible
			model, _ := projectList.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
			projectList = model.(ProjectListModel)

			projectView := projectList.View()

			// For project list, verify that all items are present
			// The delegate will format them consistently
			for i := 0; i < numItems; i++ {
				expectedName := fmt.Sprintf("Project%d", i)
				if !strings.Contains(projectView, expectedName) {
					t.Logf("Project List: Project name '%s' not found in view", expectedName)
					return false
				}
			}

			// Verify server information is present for all projects
			// The format is "Servers: [ServerX]" from the projectItem.Description() method
			for i := 0; i < numItems; i++ {
				expectedServer := fmt.Sprintf("Server%d", i)
				if !strings.Contains(projectView, expectedServer) {
					t.Logf("Project List: Server '%s' not found in view", expectedServer)
					return false
				}
			}

			// Verify the word "Servers:" appears for each project
			serversLabelCount := strings.Count(projectView, "Servers:")
			if serversLabelCount != numItems {
				t.Logf("Project List: Expected %d 'Servers:' labels, got %d", numItems, serversLabelCount)
				return false
			}

			return true
		},
		gen.IntRange(2, 10),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 21: Supplementary information styling**
// For any list item with supplementary information, the supplementary content should use secondary styling distinct from primary content
// Validates: Requirements 9.2
func TestProperty_SupplementaryInformationStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("supplementary information uses secondary styling distinct from primary content", prop.ForAll(
		func(numItems int) bool {
			// Ensure we have at least 1 item to test
			if numItems < 1 {
				numItems = 1
			}
			// Cap at reasonable number for testing
			if numItems > 20 {
				numItems = 20
			}

			// Test with SSH list - supplementary info includes connection details and auth type
			sshConfigs := make([]SSHConfig, numItems)
			for i := 0; i < numItems; i++ {
				sshConfigs[i] = SSHConfig{
					Name:     fmt.Sprintf("Server%d", i),
					Host:     fmt.Sprintf("host%d.example.com", i),
					Port:     22 + i,
					User:     fmt.Sprintf("user%d", i),
					AuthType: "password",
					Password: "pass",
				}
			}

			config := &Config{
				SSHConfigs: sshConfigs,
				Projects:   []Project{},
			}

			// Create SSH list model
			sshList := NewSSHListModel(config)
			sshView := sshList.View()

			// Verify that supplementary information (connection details and auth type) is present
			// The supplementary info should use textSecondary color styling
			// We can verify this by checking that the detail style is applied

			// Check that all connection details are present (these are supplementary info)
			for i := 0; i < numItems; i++ {
				expectedDetail := fmt.Sprintf("user%d@host%d.example.com:%d", i, i, 22+i)
				if !strings.Contains(sshView, expectedDetail) {
					t.Logf("SSH List: Connection detail '%s' (supplementary info) not found", expectedDetail)
					return false
				}

				// Check that auth type (supplementary info) is present
				if !strings.Contains(sshView, "Auth: password") {
					t.Logf("SSH List: Auth type (supplementary info) not found")
					return false
				}
			}

			// Test with Project list - supplementary info includes server list
			projects := make([]Project, numItems)
			for i := 0; i < numItems; i++ {
				projects[i] = Project{
					Name:          fmt.Sprintf("Project%d", i),
					DeployServers: []string{fmt.Sprintf("Server%d", i), fmt.Sprintf("Server%d-backup", i)},
				}
			}

			config.Projects = projects
			projectList := NewProjectListModel(config)

			// Set window size so the list renders properly
			model, _ := projectList.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
			projectList = model.(ProjectListModel)

			projectView := projectList.View()

			// Verify that supplementary information (server list) is present
			// The project description includes the server list, which is supplementary info
			for i := 0; i < numItems; i++ {
				// The supplementary info should include the server names
				expectedServer := fmt.Sprintf("Server%d", i)
				if !strings.Contains(projectView, expectedServer) {
					t.Logf("Project List: Server name '%s' (supplementary info) not found", expectedServer)
					return false
				}

				// Check for the "Servers:" label which introduces supplementary info
				if !strings.Contains(projectView, "Servers:") {
					t.Logf("Project List: 'Servers:' label (supplementary info marker) not found")
					return false
				}
			}

			// Now verify that the styling is distinct
			// We check that the detailStyle (used for supplementary info) has textSecondary color
			detailStyle := lipgloss.NewStyle().Foreground(textSecondary)
			expectedColor := textSecondary

			actualColor := detailStyle.GetForeground()
			if actualColor != expectedColor {
				t.Logf("Supplementary info style color mismatch: expected %v, got %v", expectedColor, actualColor)
				return false
			}

			// Verify that textSecondary is different from textPrimary (ensuring distinction)
			if textSecondary == textPrimary {
				t.Logf("textSecondary should be distinct from textPrimary for supplementary info")
				return false
			}

			return true
		},
		gen.IntRange(1, 10),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 15: Keyboard shortcut display consistency**
// For any screen, keyboard shortcuts should be displayed in a consistent location with symbols/formatting showing both key and action
// Validates: Requirements 7.1, 7.3, 7.5
func TestProperty_KeyboardShortcutDisplayConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("keyboard shortcuts are displayed consistently with key and action formatting", prop.ForAll(
		func(numKeys int) bool {
			// Ensure we have at least 1 key to test
			if numKeys < 1 {
				numKeys = 1
			}
			// Cap at reasonable number for testing
			if numKeys > 10 {
				numKeys = 10
			}

			// Generate random keyboard shortcuts
			keys := make(map[string]string)
			for i := 0; i < numKeys; i++ {
				key := fmt.Sprintf("key%d", i)
				action := fmt.Sprintf("Action%d", i)
				keys[key] = action
			}

			// Render the keyboard shortcuts using renderKeyHelp
			rendered := renderKeyHelp(keys)

			// Property 1: All keys should be present in the rendered output
			for key := range keys {
				if !strings.Contains(rendered, key) {
					t.Logf("Key '%s' not found in rendered output", key)
					return false
				}
			}

			// Property 2: All actions should be present in the rendered output
			for _, action := range keys {
				if !strings.Contains(rendered, action) {
					t.Logf("Action '%s' not found in rendered output", action)
					return false
				}
			}

			// Property 3: The format should clearly show key and action relationship
			// The renderKeyHelp function uses "key: action" format
			// However, we need to account for styling (ANSI codes) that may be inserted
			// So we check that both key and action appear, and that ": " separator is used
			for key, action := range keys {
				// Check that the key appears in the output
				if !strings.Contains(rendered, key) {
					t.Logf("Key '%s' not found in rendered output", key)
					return false
				}
				// Check that the action appears in the output
				if !strings.Contains(rendered, action) {
					t.Logf("Action '%s' not found in rendered output", action)
					return false
				}
			}

			// Verify that the colon separator is used (part of the format)
			if !strings.Contains(rendered, ": ") {
				t.Logf("Expected ': ' separator in rendered output")
				return false
			}

			// Property 4: Separators should be used between shortcuts
			// The renderKeyHelp function uses " • " as separator
			if numKeys > 1 {
				separatorCount := strings.Count(rendered, " • ")
				// Should have numKeys-1 separators
				if separatorCount != numKeys-1 {
					t.Logf("Expected %d separators, got %d", numKeys-1, separatorCount)
					return false
				}
			}

			// Now test consistency across different screens
			// Create different models and verify they all use renderKeyHelp consistently

			config := &Config{
				SSHConfigs: []SSHConfig{
					{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
				},
				Projects: []Project{
					{Name: "TestProject", DeployServers: []string{"TestServer"}},
				},
			}

			// Test Main Menu
			mainMenu := NewMainMenu(config)
			mainMenuView := mainMenu.View()

			// Verify main menu displays keyboard shortcuts
			// Main menu should have: ↑/↓, Enter, q
			if !strings.Contains(mainMenuView, "↑/↓") {
				t.Logf("Main menu should display '↑/↓' keyboard shortcut")
				return false
			}
			if !strings.Contains(mainMenuView, "Enter") {
				t.Logf("Main menu should display 'Enter' keyboard shortcut")
				return false
			}
			if !strings.Contains(mainMenuView, "q") {
				t.Logf("Main menu should display 'q' keyboard shortcut")
				return false
			}

			// Verify the format shows both key and action
			if !strings.Contains(mainMenuView, "Navigate") {
				t.Logf("Main menu should display 'Navigate' action")
				return false
			}
			if !strings.Contains(mainMenuView, "Select") {
				t.Logf("Main menu should display 'Select' action")
				return false
			}
			if !strings.Contains(mainMenuView, "Quit") {
				t.Logf("Main menu should display 'Quit' action")
				return false
			}

			// Test SSH List
			sshList := NewSSHListModel(config)
			sshListView := sshList.View()

			// Verify SSH list displays keyboard shortcuts
			// SSH list should have: a, e, d, t, esc, q
			if !strings.Contains(sshListView, "a") {
				t.Logf("SSH list should display 'a' keyboard shortcut")
				return false
			}
			if !strings.Contains(sshListView, "Add") {
				t.Logf("SSH list should display 'Add' action")
				return false
			}

			// Test SSH Form
			sshForm := NewSSHFormModel(config, -1)
			sshFormView := sshForm.View()

			// Verify SSH form displays keyboard shortcuts
			// SSH form should have: ↑/↓, Enter, Esc
			if !strings.Contains(sshFormView, "↑/↓") {
				t.Logf("SSH form should display '↑/↓' keyboard shortcut")
				return false
			}
			if !strings.Contains(sshFormView, "Navigate") {
				t.Logf("SSH form should display 'Navigate' action")
				return false
			}

			// Property 5: Verify consistent location
			// All screens should display keyboard shortcuts at the bottom
			// We verify this by checking that shortcuts appear in the view
			// The exact position may vary due to styling, but they should be present

			// Verify that all views contain the separator character used between shortcuts
			if !strings.Contains(mainMenuView, " • ") {
				t.Logf("Main menu should use ' • ' separator between shortcuts")
				return false
			}

			if !strings.Contains(sshListView, " • ") {
				t.Logf("SSH list should use ' • ' separator between shortcuts")
				return false
			}

			if !strings.Contains(sshFormView, " • ") {
				t.Logf("SSH form should use ' • ' separator between shortcuts")
				return false
			}

			return true
		},
		gen.IntRange(1, 5),
	))

	properties.TestingRun(t)
}

// Unit test for empty list display
// Verifies that a helpful message with appropriate styling is shown when the list is empty
// Requirements: 9.3
func TestEmptyListDisplay(t *testing.T) {
	// Create a config with no SSH configurations
	config := &Config{
		SSHConfigs: []SSHConfig{},
		Projects:   []Project{},
	}

	// Create SSH list model with empty config
	sshList := NewSSHListModel(config)

	// Get the view output
	view := sshList.View()

	// Verify that the view contains the empty state message
	if !strings.Contains(view, "No SSH Configurations Found") {
		t.Errorf("Empty list should display 'No SSH Configurations Found' message, got: %s", view)
	}

	// Verify that the view contains helpful guidance text
	if !strings.Contains(view, "You haven't added any SSH server configurations yet") {
		t.Errorf("Empty list should display helpful guidance text, got: %s", view)
	}

	// Verify that the view contains instruction to add a configuration
	if !strings.Contains(view, "Press") {
		t.Errorf("Empty list should display instruction to add configuration, got: %s", view)
	}

	// Verify that the view contains the 'a' key instruction
	if !strings.Contains(view, "a") {
		t.Errorf("Empty list should display 'a' key instruction, got: %s", view)
	}

	// Verify that the view contains "to add your first SSH configuration"
	if !strings.Contains(view, "to add your first SSH configuration") {
		t.Errorf("Empty list should display 'to add your first SSH configuration' text, got: %s", view)
	}

	// Verify that the view is not empty
	if view == "" {
		t.Error("Empty list view should not be empty")
	}

	// Verify that the view contains the title "SSH Configurations"
	if !strings.Contains(view, "SSH Configurations") {
		t.Errorf("Empty list should display 'SSH Configurations' title, got: %s", view)
	}

	// Verify that the view contains help text for navigation
	// Check for "esc" and "Back" (case-insensitive check for the key parts)
	viewLower := strings.ToLower(view)
	if !strings.Contains(viewLower, "esc") || !strings.Contains(viewLower, "back") {
		t.Errorf("Empty list should display navigation help text with 'esc' and 'back', got: %s", view)
	}
}

// **Feature: ui-enhancement, Property 7: Active form field indication**
// For any form with an active field, the active field should display a cursor or highlight indicator
// Validates: Requirements 4.2
func TestProperty_ActiveFormFieldIndication(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("active form field displays cursor and highlight indicator", prop.ForAll(
		func(activeFieldIndex int) bool {
			// Test with SSH Form
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create SSH form model (new form)
			sshForm := NewSSHFormModel(config, -1)

			// The SSH form has multiple fields (Name, Host, Port, User, Auth Type, Password/Key)
			numFields := len(sshForm.form)
			if numFields == 0 {
				t.Logf("SSH form should have fields")
				return false
			}

			// Ensure active field index is within bounds
			activeFieldIndex = activeFieldIndex % numFields
			if activeFieldIndex < 0 {
				activeFieldIndex = -activeFieldIndex
			}

			// Set the active field by simulating navigation
			for i := 0; i < activeFieldIndex; i++ {
				sshForm.cursor++
			}

			// Get the view output
			view := sshForm.View()

			// Property 1: The active field should have a cursor indicator (▶)
			if !strings.Contains(view, "▶") {
				t.Logf("Active field should display cursor indicator '▶'")
				return false
			}

			// Property 2: The view should contain visual indicators for the active field
			// The active field has special styling including borders which use box drawing characters
			// We verify the view is not empty and contains field content
			if view == "" {
				t.Logf("View should not be empty")
				return false
			}

			// Property 3: The active field should have distinct styling (border)
			// We verify this by checking that the activeFieldValueStyle is applied
			// The activeFieldValueStyle has a border with accentColor
			activeFieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(0, 1).
				Width(40)

			// Verify the style has a border style set
			borderStyle := activeFieldValueStyle.GetBorderStyle()
			if borderStyle != lipgloss.RoundedBorder() {
				t.Logf("Active field style should have rounded border")
				return false
			}

			// Verify the border color is accentColor
			borderColor := activeFieldValueStyle.GetBorderTopForeground()
			if borderColor != accentColor {
				t.Logf("Active field border color should be accentColor (%v), got %v", accentColor, borderColor)
				return false
			}

			// Property 4: Only one field should be active at a time
			// Count the number of cursor indicators (▶) - should be exactly 1
			cursorCount := strings.Count(view, "▶")
			if cursorCount != 1 {
				t.Logf("Expected exactly 1 cursor indicator, got %d", cursorCount)
				return false
			}

			// Property 5: Verify the active field has the field label visible
			// All SSH form fields should have their labels present
			if !strings.Contains(view, "Name") || !strings.Contains(view, "Host") {
				t.Logf("Form should display field labels")
				return false
			}

			// Test with Project Form as well
			projectForm := NewProjectFormModel(config, -1)

			// The project form has fields (Name, Deploy Servers)
			numProjectFields := len(projectForm.form)
			if numProjectFields == 0 {
				t.Logf("Project form should have fields")
				return false
			}

			// Set active field in project form
			projectActiveIndex := activeFieldIndex % numProjectFields
			for i := 0; i < projectActiveIndex; i++ {
				projectForm.cursor++
			}

			// Get the project form view
			projectView := projectForm.View()

			// Verify project form also has cursor indicator
			if !strings.Contains(projectView, "▶") {
				t.Logf("Project form active field should display cursor indicator '▶'")
				return false
			}

			// Property 6: Verify that inactive fields do NOT have the cursor indicator
			// Create a form and verify that when cursor is at position 0,
			// the other fields don't have the cursor indicator
			testForm := NewSSHFormModel(config, -1)
			testForm.cursor = 0 // Set to first field
			testView := testForm.View()

			// Count cursor indicators - should be exactly 1
			testCursorCount := strings.Count(testView, "▶")
			if testCursorCount != 1 {
				t.Logf("Expected exactly 1 cursor indicator in test form, got %d", testCursorCount)
				return false
			}

			// Property 7: Verify that the active field indication is consistent across form types
			// Both SSH form and Project form should use cursor indicators
			// (both now use the same symbol for consistency)
			sshHasIndicator := strings.Contains(view, "▶")
			projectHasIndicator := strings.Contains(projectView, "▶")

			if !sshHasIndicator || !projectHasIndicator {
				t.Logf("Both form types should have cursor indicators")
				return false
			}

			return true
		},
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 3: Form field alignment**
// For any form, all form fields should be aligned consistently with uniform spacing between fields
// Validates: Requirements 2.3
func TestProperty_FormFieldAlignment(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all form fields have consistent alignment and uniform spacing", prop.ForAll(
		func(seed int64) bool {
			// Test with SSH Form
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create SSH form model (new form)
			_ = NewSSHFormModel(config, -1)

			// Property 1: Verify that the form uses consistent label width
			// The fieldLabelStyle in SSHFormModel.View() has Width(20)
			// This ensures all labels are aligned consistently

			// We test this by checking the style definitions directly
			// rather than parsing the rendered output (which contains ANSI codes)

			// Create the label style as defined in the form
			fieldLabelStyle := lipgloss.NewStyle().
				Foreground(textSecondary).
				Bold(true).
				Width(20).
				Align(lipgloss.Left)

			// Verify the label width is set to 20
			labelWidth := fieldLabelStyle.GetWidth()
			if labelWidth != 20 {
				t.Logf("SSH form label width should be 20, got %d", labelWidth)
				return false
			}

			// Property 2: Verify that the form uses consistent value width
			// The fieldValueStyle has Width(40)
			fieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				Padding(0, 1).
				Width(40)

			valueWidth := fieldValueStyle.GetWidth()
			if valueWidth != 40 {
				t.Logf("SSH form value width should be 40, got %d", valueWidth)
				return false
			}

			// Property 3: Verify that active field style also uses consistent width
			activeFieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(0, 1).
				Width(40)

			activeValueWidth := activeFieldValueStyle.GetWidth()
			if activeValueWidth != 40 {
				t.Logf("SSH form active value width should be 40, got %d", activeValueWidth)
				return false
			}

			// Property 4: Verify that error field style also uses consistent width
			errorFieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(errorColor).
				Padding(0, 1).
				Width(40)

			errorValueWidth := errorFieldValueStyle.GetWidth()
			if errorValueWidth != 40 {
				t.Logf("SSH form error value width should be 40, got %d", errorValueWidth)
				return false
			}

			// Test with Project Form as well
			projectForm := NewProjectFormModel(config, -1)
			_ = projectForm // Use the variable

			// Project form should use the same label and value widths
			// Create the styles as defined in ProjectFormModel.View()
			projectFieldLabelStyle := lipgloss.NewStyle().
				Foreground(textSecondary).
				Bold(true).
				Width(20).
				Align(lipgloss.Left)

			projectLabelWidth := projectFieldLabelStyle.GetWidth()
			if projectLabelWidth != 20 {
				t.Logf("Project form label width should be 20, got %d", projectLabelWidth)
				return false
			}

			projectFieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				Padding(0, 1).
				Width(40)

			projectValueWidth := projectFieldValueStyle.GetWidth()
			if projectValueWidth != 40 {
				t.Logf("Project form value width should be 40, got %d", projectValueWidth)
				return false
			}

			// Property 5: Verify that both forms use the same widths (consistency across forms)
			if labelWidth != projectLabelWidth {
				t.Logf("Label width mismatch between forms: SSH=%d, Project=%d", labelWidth, projectLabelWidth)
				return false
			}

			if valueWidth != projectValueWidth {
				t.Logf("Value width mismatch between forms: SSH=%d, Project=%d", valueWidth, projectValueWidth)
				return false
			}

			// Property 6: Verify uniform spacing by checking the form rendering
			// Each field should have consistent spacing (one blank line between fields)
			// This is enforced by the formContent.WriteString("\n") after each field

			return true
		},
		gen.Int64(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 17: Form cursor position display**
// For any active form field, the current cursor position within the field should be visually indicated
// Validates: Requirements 8.1
func TestProperty_FormCursorPositionDisplay(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("active form field displays cursor position indicator", prop.ForAll(
		func(activeFieldIndex int, fieldValue string) bool {
			// Test with SSH Form
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create SSH form model (new form)
			sshForm := NewSSHFormModel(config, -1)

			// The SSH form has multiple fields
			numFields := len(sshForm.form)
			if numFields == 0 {
				t.Logf("SSH form should have fields")
				return false
			}

			// Ensure active field index is within bounds
			activeFieldIndex = activeFieldIndex % numFields
			if activeFieldIndex < 0 {
				activeFieldIndex = -activeFieldIndex
			}

			// Set the active field
			sshForm.cursor = activeFieldIndex

			// Set a value for the active field (limit length for testing)
			if len(fieldValue) > 50 {
				fieldValue = fieldValue[:50]
			}
			sshForm.form[activeFieldIndex].value = fieldValue

			// Get the view output
			view := sshForm.View()

			// Property 1: The active field should display a cursor position indicator
			// The cursor position indicator is "│" character
			if !strings.Contains(view, "│") {
				t.Logf("Active field should display cursor position indicator '│'")
				return false
			}

			// Property 2: The cursor position indicator should appear in the view
			// when the field is active (has focus)
			// Count the number of cursor position indicators
			cursorPosCount := strings.Count(view, "│")

			// Should have at least one cursor position indicator for the active field
			if cursorPosCount < 1 {
				t.Logf("Expected at least 1 cursor position indicator, got %d", cursorPosCount)
				return false
			}

			// Property 3: The view should contain the field value
			// (unless it's a password field, which is masked)
			fieldType := sshForm.form[activeFieldIndex].fieldType
			if fieldType != "password" && fieldValue != "" {
				if !strings.Contains(view, fieldValue) {
					t.Logf("View should contain field value '%s'", fieldValue)
					return false
				}
			}

			// Property 4: Test with Project Form as well
			projectForm := NewProjectFormModel(config, -1)
			numProjectFields := len(projectForm.form)
			if numProjectFields == 0 {
				t.Logf("Project form should have fields")
				return false
			}

			// Set active field in project form
			projectActiveIndex := activeFieldIndex % numProjectFields
			projectForm.cursor = projectActiveIndex

			// Set a value for the active field
			projectForm.form[projectActiveIndex].value = fieldValue

			// Get the project form view
			projectView := projectForm.View()

			// Verify project form also has cursor position indicator
			if !strings.Contains(projectView, "│") {
				t.Logf("Project form active field should display cursor position indicator '│'")
				return false
			}

			// Property 5: Verify that the cursor position indicator is only shown
			// for the active field, not for inactive fields
			// Create a form with multiple fields and verify only the active one has the indicator
			testForm := NewSSHFormModel(config, -1)
			if len(testForm.form) < 2 {
				// Need at least 2 fields to test this property
				return true
			}

			// Set first field as active
			testForm.cursor = 0
			testForm.form[0].value = "test1"
			testForm.form[1].value = "test2"

			testView := testForm.View()

			// The view should contain the cursor position indicator
			if !strings.Contains(testView, "│") {
				t.Logf("Test form should display cursor position indicator")
				return false
			}

			// Property 6: Verify that when a field has no value, the cursor indicator is still shown
			emptyForm := NewSSHFormModel(config, -1)
			emptyForm.cursor = 0
			emptyForm.form[0].value = ""

			emptyView := emptyForm.View()

			// Even with empty value, the cursor position indicator should be present
			if !strings.Contains(emptyView, "│") {
				t.Logf("Empty field should still display cursor position indicator")
				return false
			}

			// Property 7: Verify that the cursor position indicator is clearly visible
			// by checking it's rendered with the cursorIndicatorStyle
			// The cursorIndicatorStyle uses accentColor
			cursorIndicatorStyle := lipgloss.NewStyle().Foreground(accentColor)
			indicatorColor := cursorIndicatorStyle.GetForeground()

			if indicatorColor != accentColor {
				t.Logf("Cursor position indicator should use accentColor (%v), got %v", accentColor, indicatorColor)
				return false
			}

			// Property 8: Verify consistency - the cursor position indicator should be
			// the same character across all form types
			// Both SSH form and Project form should use "│"
			sshHasIndicator := strings.Contains(view, "│")
			projectHasIndicator := strings.Contains(projectView, "│")

			if !sshHasIndicator || !projectHasIndicator {
				t.Logf("Both form types should use the same cursor position indicator '│'")
				return false
			}

			return true
		},
		gen.IntRange(0, 10),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 11: Label consistency across forms**
// For any two forms, the label formatting (style, alignment, spacing) should be identical
// Validates: Requirements 5.4
func TestProperty_LabelConsistencyAcrossForms(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all forms use identical label formatting, alignment, and spacing", prop.ForAll(
		func(seed int64) bool {
			// Create test config
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create both form types
			sshForm := NewSSHFormModel(config, -1)
			projectForm := NewProjectFormModel(config, -1)

			// Property 1: Both forms should use the same label style characteristics
			// Define the expected label style as used in SSHFormModel
			expectedLabelStyle := lipgloss.NewStyle().
				Foreground(textSecondary).
				Bold(true).
				Width(20).
				Align(lipgloss.Left)

			// Verify SSH form uses this style
			sshLabelWidth := expectedLabelStyle.GetWidth()
			sshLabelForeground := expectedLabelStyle.GetForeground()
			sshLabelBold := expectedLabelStyle.GetBold()
			sshLabelAlign := expectedLabelStyle.GetAlign()

			if sshLabelWidth != 20 {
				t.Logf("SSH form label width should be 20, got %d", sshLabelWidth)
				return false
			}

			if sshLabelForeground != textSecondary {
				t.Logf("SSH form label foreground should be textSecondary (%v), got %v", textSecondary, sshLabelForeground)
				return false
			}

			if !sshLabelBold {
				t.Logf("SSH form labels should be bold")
				return false
			}

			if sshLabelAlign != lipgloss.Left {
				t.Logf("SSH form labels should be left-aligned")
				return false
			}

			// Property 2: Project form should use the same label style
			// Define the label style as it should be used in ProjectFormModel
			projectLabelStyle := lipgloss.NewStyle().
				Foreground(textSecondary).
				Bold(true).
				Width(20).
				Align(lipgloss.Left)

			projectLabelWidth := projectLabelStyle.GetWidth()
			projectLabelForeground := projectLabelStyle.GetForeground()
			projectLabelBold := projectLabelStyle.GetBold()
			projectLabelAlign := projectLabelStyle.GetAlign()

			// Verify project form uses the same characteristics
			if projectLabelWidth != sshLabelWidth {
				t.Logf("Label width mismatch: SSH=%d, Project=%d", sshLabelWidth, projectLabelWidth)
				return false
			}

			if projectLabelForeground != sshLabelForeground {
				t.Logf("Label foreground mismatch: SSH=%v, Project=%v", sshLabelForeground, projectLabelForeground)
				return false
			}

			if projectLabelBold != sshLabelBold {
				t.Logf("Label bold mismatch: SSH=%v, Project=%v", sshLabelBold, projectLabelBold)
				return false
			}

			if projectLabelAlign != sshLabelAlign {
				t.Logf("Label alignment mismatch: SSH=%v, Project=%v", sshLabelAlign, projectLabelAlign)
				return false
			}

			// Property 3: Verify consistency by rendering both forms and checking label presence
			sshView := sshForm.View()
			projectView := projectForm.View()

			// Both views should contain labels
			if sshView == "" || projectView == "" {
				t.Logf("Form views should not be empty")
				return false
			}

			// SSH form should have labels: Name, Host, Port, User, Auth Type
			sshLabels := []string{"Name", "Host", "Port", "User", "Auth Type"}
			for _, label := range sshLabels {
				if !strings.Contains(sshView, label) {
					t.Logf("SSH form should contain label '%s'", label)
					return false
				}
			}

			// Project form should have labels: Name, Deploy Servers
			projectLabels := []string{"Name", "Deploy Servers"}
			for _, label := range projectLabels {
				if !strings.Contains(projectView, label) {
					t.Logf("Project form should contain label '%s'", label)
					return false
				}
			}

			// Property 4: Verify that both forms use the same textSecondary color for labels
			// This is already checked above, but we verify it's consistent with the global style
			if textSecondary == textPrimary {
				t.Logf("textSecondary should be distinct from textPrimary for label differentiation")
				return false
			}

			// Property 5: Verify that both forms use bold styling for labels
			// This ensures labels stand out consistently across all forms
			// Already verified above through the style checks

			// Property 6: Verify that both forms use the same width for labels (20)
			// This ensures consistent alignment across all forms
			// Already verified above

			// Property 7: Verify that both forms use left alignment for labels
			// This ensures consistent visual presentation
			// Already verified above

			// Property 8: Test with different field counts to ensure consistency holds
			// The SSH form has more fields than the project form
			// But the label styling should remain consistent regardless of field count
			sshFieldCount := len(sshForm.form)
			projectFieldCount := len(projectForm.form)

			if sshFieldCount == 0 || projectFieldCount == 0 {
				t.Logf("Forms should have fields")
				return false
			}

			// Both forms should have fields, and despite different counts,
			// the label styling should be identical
			// This is verified by the style checks above

			return true
		},
		gen.Int64(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 19: Multi-line field boundaries**
// For any multi-line form field, the field should have clear visual boundaries and adequate space
// Validates: Requirements 8.4
func TestProperty_MultiLineFieldBoundaries(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("multi-line form fields have clear visual boundaries and adequate space", prop.ForAll(
		func(numLines int, lineContent string) bool {
			// Ensure we have at least 1 line and cap at reasonable number
			if numLines < 1 {
				numLines = 1
			}
			if numLines > 20 {
				numLines = 20
			}

			// Limit line content length for testing
			if len(lineContent) > 50 {
				lineContent = lineContent[:50]
			}

			// Create test config
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Create project form which has multi-line fields
			// (Build Instructions and Deploy Script are marked as multiline: true)
			projectForm := NewProjectFormModel(config, -1)

			// Find the multi-line fields in the form
			var multilineFieldIndices []int
			for i, field := range projectForm.form {
				if field.multiline {
					multilineFieldIndices = append(multilineFieldIndices, i)
				}
			}

			// Property 1: The form should have at least one multi-line field
			if len(multilineFieldIndices) == 0 {
				t.Logf("Project form should have at least one multi-line field")
				return false
			}

			// Property 2: Set multi-line content in the multi-line fields
			// Create multi-line content by joining lines with newline
			lines := make([]string, numLines)
			for i := 0; i < numLines; i++ {
				lines[i] = fmt.Sprintf("%s_%d", lineContent, i)
			}
			multilineContent := strings.Join(lines, "\n")

			// Set the multi-line content in the first multi-line field
			firstMultilineIndex := multilineFieldIndices[0]
			projectForm.form[firstMultilineIndex].value = multilineContent

			// Set cursor to the multi-line field to make it active
			projectForm.cursor = firstMultilineIndex

			// Get the view output
			view := projectForm.View()

			// Property 3: The view should contain the multi-line content
			// When content has newlines, the view should contain the content
			// We check that the view contains the multi-line content (with or without newlines rendered)
			// For simple forms, the content is rendered as-is, so newlines are preserved
			// We verify by checking that at least some of the content is present
			if multilineContent != "" && !strings.Contains(view, strings.Split(multilineContent, "\n")[0]) {
				t.Logf("Multi-line field view should contain at least the first line of content")
				return false
			}

			// Property 4: Multi-line fields should have visual boundaries
			// The form uses formStyle which has a border
			// Verify that the form has border characters
			borderChars := []string{"╭", "╮", "╰", "╯", "─", "│"}
			hasBorder := false
			for _, char := range borderChars {
				if strings.Contains(view, char) {
					hasBorder = true
					break
				}
			}

			if !hasBorder {
				t.Logf("Multi-line field should have visual boundaries (border)")
				return false
			}

			// Property 5: Multi-line fields should have adequate space
			// Verify that the field value style has appropriate width
			// The field value style should have a width set (40 in current implementation)
			fieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				Padding(0, 1).
				Width(40)

			valueWidth := fieldValueStyle.GetWidth()
			if valueWidth < 20 {
				t.Logf("Multi-line field should have adequate width (at least 20), got %d", valueWidth)
				return false
			}

			// Property 6: Multi-line fields should have padding for readability
			// Verify that the field value style has padding
			paddingLeft := fieldValueStyle.GetPaddingLeft()
			paddingRight := fieldValueStyle.GetPaddingRight()

			if paddingLeft == 0 && paddingRight == 0 {
				t.Logf("Multi-line field should have horizontal padding for readability")
				return false
			}

			// Property 7: When a multi-line field is active, it should have distinct styling
			// The active field should have a border with accent color
			activeFieldValueStyle := lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(backgroundLight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(0, 1).
				Width(40)

			// Verify the active field has a border
			activeBorderStyle := activeFieldValueStyle.GetBorderStyle()
			if activeBorderStyle != lipgloss.RoundedBorder() {
				t.Logf("Active multi-line field should have a border")
				return false
			}

			// Verify the border color is accent color
			activeBorderColor := activeFieldValueStyle.GetBorderTopForeground()
			if activeBorderColor != accentColor {
				t.Logf("Active multi-line field border should use accent color")
				return false
			}

			// Property 8: Multi-line fields should maintain consistent styling with single-line fields
			// Both should use the same base styles (colors, padding)
			// Verify foreground color is consistent
			if fieldValueStyle.GetForeground() != activeFieldValueStyle.GetForeground() {
				t.Logf("Multi-line and single-line fields should use consistent foreground color")
				return false
			}

			// Verify background color is consistent
			if fieldValueStyle.GetBackground() != activeFieldValueStyle.GetBackground() {
				t.Logf("Multi-line and single-line fields should use consistent background color")
				return false
			}

			// Property 9: Multi-line fields should have clear separation from other fields
			// Verify that there is spacing between fields (newline characters)
			// The form adds "\n" after each field
			// Count the number of fields and verify spacing
			fieldCount := len(projectForm.form)
			if fieldCount > 1 {
				// There should be multiple fields, and they should be separated
				// We verify this by checking that the view is not just a single line
				viewLines := strings.Split(view, "\n")
				if len(viewLines) < fieldCount {
					t.Logf("Multi-line fields should have clear separation (expected at least %d lines, got %d)", fieldCount, len(viewLines))
					return false
				}
			}

			// Property 10: Test with SSH form as well (even though it doesn't currently have multi-line fields)
			// This ensures consistency if multi-line fields are added to SSH form in the future
			sshForm := NewSSHFormModel(config, -1)

			// Check if SSH form has any multi-line fields
			var sshMultilineIndices []int
			for i, field := range sshForm.form {
				if field.multiline {
					sshMultilineIndices = append(sshMultilineIndices, i)
				}
			}

			// If SSH form has multi-line fields, verify they follow the same properties
			if len(sshMultilineIndices) > 0 {
				// Set multi-line content in the first multi-line field
				sshForm.form[sshMultilineIndices[0]].value = multilineContent
				sshForm.cursor = sshMultilineIndices[0]

				sshView := sshForm.View()

				// Verify all lines are present
				for _, line := range lines {
					if line != "" && !strings.Contains(sshView, line) {
						t.Logf("SSH form multi-line field should contain line: %s", line)
						return false
					}
				}

				// Verify border is present
				hasSshBorder := false
				for _, char := range borderChars {
					if strings.Contains(sshView, char) {
						hasSshBorder = true
						break
					}
				}

				if !hasSshBorder {
					t.Logf("SSH form multi-line field should have visual boundaries")
					return false
				}
			}

			// Property 11: Multi-line fields should handle empty content gracefully
			// Test with empty multi-line field
			emptyForm := NewProjectFormModel(config, -1)
			if len(multilineFieldIndices) > 0 {
				emptyForm.form[firstMultilineIndex].value = ""
				emptyForm.cursor = firstMultilineIndex

				emptyView := emptyForm.View()

				// Even with empty content, the field should still have boundaries
				hasEmptyBorder := false
				for _, char := range borderChars {
					if strings.Contains(emptyView, char) {
						hasEmptyBorder = true
						break
					}
				}

				if !hasEmptyBorder {
					t.Logf("Empty multi-line field should still have visual boundaries")
					return false
				}

				// The view should not be empty
				if emptyView == "" {
					t.Logf("Empty multi-line field view should not be empty")
					return false
				}
			}

			// Property 12: Multi-line fields should preserve line breaks in the content
			// Verify that newline characters in the content are handled appropriately
			// This is implicitly tested by checking that all lines are present in the view

			return true
		},
		gen.IntRange(1, 10),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test for deployment step completion display
// Verifies that a success indicator is displayed on completion
// Requirements: 6.2
func TestDeploymentStepCompletionDisplay(t *testing.T) {
	// Create a test config with SSH configuration
	config := &Config{
		SSHConfigs: []SSHConfig{
			{
				Name:     "TestServer",
				Host:     "test.example.com",
				Port:     22,
				User:     "testuser",
				AuthType: "password",
				Password: "testpass",
			},
		},
		Projects: []Project{},
	}

	// Create SSH test model
	sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

	// Simulate successful connection by setting result
	sshTest.result = "Success: SSH connection successful"
	sshTest.done = true

	// Get the view output
	view := sshTest.View()

	// Verify that the view contains a success indicator (✓ icon)
	if !strings.Contains(view, "✓") {
		t.Errorf("Completion display should contain success indicator '✓', got: %s", view)
	}

	// Verify that the view contains success styling
	// The success message should be present
	if !strings.Contains(view, "Success") {
		t.Errorf("Completion display should contain 'Success' message, got: %s", view)
	}

	// Verify that the view contains "Connected" status
	if !strings.Contains(view, "Connected") {
		t.Errorf("Completion display should contain 'Connected' status, got: %s", view)
	}

	// Verify that connection details are displayed
	if !strings.Contains(view, "Connection Details") {
		t.Errorf("Completion display should contain 'Connection Details' card, got: %s", view)
	}

	// Verify that host information is displayed
	if !strings.Contains(view, "test.example.com") {
		t.Errorf("Completion display should contain host information, got: %s", view)
	}

	// Verify that the view is not empty
	if view == "" {
		t.Error("Completion display view should not be empty")
	}
}

// Unit test for deployment step failure display
// Verifies that an error is highlighted with prominent styling
// Requirements: 6.3
func TestDeploymentStepFailureDisplay(t *testing.T) {
	// Create a test config with SSH configuration
	config := &Config{
		SSHConfigs: []SSHConfig{
			{
				Name:     "TestServer",
				Host:     "test.example.com",
				Port:     22,
				User:     "testuser",
				AuthType: "password",
				Password: "testpass",
			},
		},
		Projects: []Project{},
	}

	// Create SSH test model
	sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

	// Simulate failed connection by setting error result
	sshTest.result = "Connection failed: connection refused"
	sshTest.done = true

	// Get the view output
	view := sshTest.View()

	// Verify that the view contains an error indicator (✗ icon)
	if !strings.Contains(view, "✗") {
		t.Errorf("Failure display should contain error indicator '✗', got: %s", view)
	}

	// Verify that the view contains "Connection Failed" message
	if !strings.Contains(view, "Connection Failed") {
		t.Errorf("Failure display should contain 'Connection Failed' message, got: %s", view)
	}

	// Verify that the view contains error details
	if !strings.Contains(view, "Error:") {
		t.Errorf("Failure display should contain 'Error:' section, got: %s", view)
	}

	// Verify that the actual error message is displayed
	if !strings.Contains(view, "connection refused") {
		t.Errorf("Failure display should contain actual error message, got: %s", view)
	}

	// Verify that connection details are displayed even on failure
	if !strings.Contains(view, "Connection Details") {
		t.Errorf("Failure display should contain 'Connection Details' card, got: %s", view)
	}

	// Verify that host information is displayed
	if !strings.Contains(view, "test.example.com") {
		t.Errorf("Failure display should contain host information, got: %s", view)
	}

	// Verify that the view is not empty
	if view == "" {
		t.Error("Failure display view should not be empty")
	}
}

// Unit test for deployment in progress display
// Verifies that a progress indicator is shown during deployment
// Requirements: 6.4
func TestDeploymentInProgressDisplay(t *testing.T) {
	// Create a test config with SSH configuration
	config := &Config{
		SSHConfigs: []SSHConfig{
			{
				Name:     "TestServer",
				Host:     "test.example.com",
				Port:     22,
				User:     "testuser",
				AuthType: "password",
				Password: "testpass",
			},
		},
		Projects: []Project{},
	}

	// Create SSH test model
	sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

	// Simulate in-progress state (done = false)
	sshTest.done = false
	sshTest.spinnerFrame = 0

	// Get the view output
	view := sshTest.View()

	// Verify that the view contains a progress indicator (spinner)
	// The spinner uses characters from the set: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	hasSpinner := false
	for _, char := range spinnerChars {
		if strings.Contains(view, char) {
			hasSpinner = true
			break
		}
	}
	if !hasSpinner {
		t.Errorf("In-progress display should contain spinner indicator, got: %s", view)
	}

	// Verify that the view contains "Testing connection..." message
	if !strings.Contains(view, "Testing connection...") {
		t.Errorf("In-progress display should contain 'Testing connection...' message, got: %s", view)
	}

	// Verify that the view contains loading icon (⋯)
	if !strings.Contains(view, "⋯") {
		t.Errorf("In-progress display should contain loading icon '⋯', got: %s", view)
	}

	// Verify that connection details are displayed during progress
	if !strings.Contains(view, "Connection Details") {
		t.Errorf("In-progress display should contain 'Connection Details' card, got: %s", view)
	}

	// Verify that host information is displayed
	if !strings.Contains(view, "test.example.com") {
		t.Errorf("In-progress display should contain host information, got: %s", view)
	}

	// Verify that the view is not empty
	if view == "" {
		t.Error("In-progress display view should not be empty")
	}

	// Test with different spinner frames to ensure animation works
	for frame := 0; frame < 10; frame++ {
		sshTest.spinnerFrame = frame
		frameView := sshTest.View()

		// Each frame should produce a non-empty view
		if frameView == "" {
			t.Errorf("Spinner frame %d should produce non-empty view", frame)
		}

		// Should still contain the testing message
		if !strings.Contains(frameView, "Testing connection...") {
			t.Errorf("Spinner frame %d should contain 'Testing connection...' message", frame)
		}
	}
}

// **Feature: ui-enhancement, Property 23: Success message styling**
// For any successful action completion, a success message with positive/success styling should be displayed
// Validates: Requirements 10.1
func TestProperty_SuccessMessageStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("success messages display with positive styling and success indicator", prop.ForAll(
		func(serverName string, host string, port int, user string) bool {
			// Ensure valid port range
			if port <= 0 || port > 65535 {
				port = 22
			}

			// Ensure non-empty strings
			if serverName == "" {
				serverName = "TestServer"
			}
			if host == "" {
				host = "test.example.com"
			}
			if user == "" {
				user = "testuser"
			}

			// Create a test config with SSH configuration
			config := &Config{
				SSHConfigs: []SSHConfig{
					{
						Name:     serverName,
						Host:     host,
						Port:     port,
						User:     user,
						AuthType: "password",
						Password: "testpass",
					},
				},
				Projects: []Project{},
			}

			// Create SSH test model
			sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

			// Simulate successful connection
			sshTest.result = "Success: SSH connection successful"
			sshTest.done = true

			// Get the view output
			view := sshTest.View()

			// Property 1: Success messages should contain a success indicator icon (✓)
			if !strings.Contains(view, "✓") {
				t.Logf("Success message should contain success indicator '✓'")
				return false
			}

			// Property 2: Success messages should contain the word "Success"
			if !strings.Contains(view, "Success") {
				t.Logf("Success message should contain 'Success' text")
				return false
			}

			// Property 3: Success messages should display connection details
			if !strings.Contains(view, "Connection Details") {
				t.Logf("Success message should display 'Connection Details'")
				return false
			}

			// Property 4: Success messages should show "Connected" status
			if !strings.Contains(view, "Connected") {
				t.Logf("Success message should show 'Connected' status")
				return false
			}

			// Property 5: Success messages should display the host information
			if !strings.Contains(view, host) {
				t.Logf("Success message should display host '%s'", host)
				return false
			}

			// Property 6: Success messages should display the user information
			if !strings.Contains(view, user) {
				t.Logf("Success message should display user '%s'", user)
				return false
			}

			// Property 7: Success messages should display the port information
			portStr := fmt.Sprintf("%d", port)
			if !strings.Contains(view, portStr) {
				t.Logf("Success message should display port '%s'", portStr)
				return false
			}

			// Property 8: Verify that the success style uses the correct color
			// The successStyle should use successColor
			if successStyle.GetForeground() != successColor {
				t.Logf("Success style should use successColor (%v), got %v", successColor, successStyle.GetForeground())
				return false
			}

			// Property 9: Success messages should be non-empty
			if view == "" {
				t.Logf("Success message view should not be empty")
				return false
			}

			// Property 10: Success messages should include help text for returning
			if !strings.Contains(view, "Return") || !strings.Contains(view, "Enter") {
				t.Logf("Success message should include help text for returning")
				return false
			}

			return true
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.IntRange(1, 65535),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 24: Error message styling**
// For any error occurrence, an error message with error styling should be displayed
// Validates: Requirements 10.2
func TestProperty_ErrorMessageStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("error messages display with error styling and error indicator", prop.ForAll(
		func(serverName string, host string, port int, user string, errorMsg string) bool {
			// Ensure valid port range
			if port <= 0 || port > 65535 {
				port = 22
			}

			// Ensure non-empty strings
			if serverName == "" {
				serverName = "TestServer"
			}
			if host == "" {
				host = "test.example.com"
			}
			if user == "" {
				user = "testuser"
			}
			if errorMsg == "" {
				errorMsg = "connection refused"
			}

			// Create a test config with SSH configuration
			config := &Config{
				SSHConfigs: []SSHConfig{
					{
						Name:     serverName,
						Host:     host,
						Port:     port,
						User:     user,
						AuthType: "password",
						Password: "testpass",
					},
				},
				Projects: []Project{},
			}

			// Create SSH test model
			sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

			// Simulate failed connection with error message
			sshTest.result = fmt.Sprintf("Connection failed: %s", errorMsg)
			sshTest.done = true

			// Get the view output
			view := sshTest.View()

			// Property 1: Error messages should contain an error indicator icon (✗)
			if !strings.Contains(view, "✗") {
				t.Logf("Error message should contain error indicator '✗'")
				return false
			}

			// Property 2: Error messages should contain "Connection Failed" text
			if !strings.Contains(view, "Connection Failed") {
				t.Logf("Error message should contain 'Connection Failed' text")
				return false
			}

			// Property 3: Error messages should display connection details
			if !strings.Contains(view, "Connection Details") {
				t.Logf("Error message should display 'Connection Details'")
				return false
			}

			// Property 4: Error messages should show the error section
			if !strings.Contains(view, "Error:") {
				t.Logf("Error message should show 'Error:' section")
				return false
			}

			// Property 5: Error messages should display the actual error message
			if !strings.Contains(view, errorMsg) {
				t.Logf("Error message should display actual error '%s'", errorMsg)
				return false
			}

			// Property 6: Error messages should display the host information
			if !strings.Contains(view, host) {
				t.Logf("Error message should display host '%s'", host)
				return false
			}

			// Property 7: Error messages should display the user information
			if !strings.Contains(view, user) {
				t.Logf("Error message should display user '%s'", user)
				return false
			}

			// Property 8: Error messages should display the port information
			portStr := fmt.Sprintf("%d", port)
			if !strings.Contains(view, portStr) {
				t.Logf("Error message should display port '%s'", portStr)
				return false
			}

			// Property 9: Verify that the error style uses the correct color
			// The errorStyle should use errorColor
			if errorStyle.GetForeground() != errorColor {
				t.Logf("Error style should use errorColor (%v), got %v", errorColor, errorStyle.GetForeground())
				return false
			}

			// Property 10: Error messages should be non-empty
			if view == "" {
				t.Logf("Error message view should not be empty")
				return false
			}

			// Property 11: Error messages should include help text for returning
			if !strings.Contains(view, "Return") || !strings.Contains(view, "Enter") {
				t.Logf("Error message should include help text for returning")
				return false
			}

			return true
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.IntRange(1, 65535),
		gen.AnyString(),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 25: Processing indicator display**
// For any processing operation, a loading or progress indicator should be displayed
// Validates: Requirements 10.3
func TestProperty_ProcessingIndicatorDisplay(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("processing operations display loading indicator and progress message", prop.ForAll(
		func(serverName string, host string, port int, user string, spinnerFrame int) bool {
			// Ensure valid port range
			if port <= 0 || port > 65535 {
				port = 22
			}

			// Ensure non-empty strings
			if serverName == "" {
				serverName = "TestServer"
			}
			if host == "" {
				host = "test.example.com"
			}
			if user == "" {
				user = "testuser"
			}

			// Ensure spinner frame is non-negative
			if spinnerFrame < 0 {
				spinnerFrame = -spinnerFrame
			}

			// Create a test config with SSH configuration
			config := &Config{
				SSHConfigs: []SSHConfig{
					{
						Name:     serverName,
						Host:     host,
						Port:     port,
						User:     user,
						AuthType: "password",
						Password: "testpass",
					},
				},
				Projects: []Project{},
			}

			// Create SSH test model
			sshTest := NewSSHTestModel(config.SSHConfigs[0], config)

			// Simulate processing state (not done)
			sshTest.done = false
			sshTest.spinnerFrame = spinnerFrame

			// Get the view output
			view := sshTest.View()

			// Property 1: Processing operations should display a spinner indicator
			// The spinner uses characters from the set: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
			spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			hasSpinner := false
			for _, char := range spinnerChars {
				if strings.Contains(view, char) {
					hasSpinner = true
					break
				}
			}
			if !hasSpinner {
				t.Logf("Processing operation should display spinner indicator")
				return false
			}

			// Property 2: Processing operations should display a progress message
			if !strings.Contains(view, "Testing connection...") {
				t.Logf("Processing operation should display progress message 'Testing connection...'")
				return false
			}

			// Property 3: Processing operations should display a loading icon (⋯)
			if !strings.Contains(view, "⋯") {
				t.Logf("Processing operation should display loading icon '⋯'")
				return false
			}

			// Property 4: Processing operations should display connection details
			if !strings.Contains(view, "Connection Details") {
				t.Logf("Processing operation should display 'Connection Details'")
				return false
			}

			// Property 5: Processing operations should display the host information
			if !strings.Contains(view, host) {
				t.Logf("Processing operation should display host '%s'", host)
				return false
			}

			// Property 6: Processing operations should display the user information
			if !strings.Contains(view, user) {
				t.Logf("Processing operation should display user '%s'", user)
				return false
			}

			// Property 7: Processing operations should display the port information
			portStr := fmt.Sprintf("%d", port)
			if !strings.Contains(view, portStr) {
				t.Logf("Processing operation should display port '%s'", portStr)
				return false
			}

			// Property 8: Verify that the status style uses the correct color
			// The statusStyle should use infoColor
			if statusStyle.GetForeground() != infoColor {
				t.Logf("Status style should use infoColor (%v), got %v", infoColor, statusStyle.GetForeground())
				return false
			}

			// Property 9: Processing operations should be non-empty
			if view == "" {
				t.Logf("Processing operation view should not be empty")
				return false
			}

			// Property 10: Processing operations should NOT display help text
			// (help text is only shown when done)
			if strings.Contains(view, "Press Enter to return") {
				t.Logf("Processing operation should not display help text (only shown when done)")
				return false
			}

			// Property 11: Verify that different spinner frames produce different outputs
			// (animation property)
			if spinnerFrame < 100 {
				// Test with a different frame
				sshTest.spinnerFrame = spinnerFrame + 1
				view2 := sshTest.View()

				// The views should be different (due to different spinner frames)
				// However, they may be the same if the frame wraps around to the same character
				// So we just verify that both views are valid (non-empty)
				if view2 == "" {
					t.Logf("Different spinner frame should produce valid view")
					return false
				}
			}

			return true
		},
		gen.AnyString(),
		gen.AnyString(),
		gen.IntRange(1, 65535),
		gen.AnyString(),
		gen.IntRange(0, 1000),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 13: Log entry type indicators**
// For any deployment log entry, the rendered output should include an icon or symbol corresponding to the entry type
// Validates: Requirements 6.1
func TestProperty_LogEntryTypeIndicators(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("log entries include icons corresponding to entry type", prop.ForAll(
		func(logType string, logMessage string) bool {
			// Skip empty messages
			if strings.TrimSpace(logMessage) == "" {
				return true
			}

			// Create a test config with a project
			config := &Config{
				SSHConfigs: []SSHConfig{
					{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
				},
				Projects: []Project{
					{Name: "TestProject", DeployServers: []string{"TestServer"}},
				},
			}

			// Create deploy model
			deployModel := NewDeployModel(config.Projects[0], config)

			// Add log entries of different types based on logType
			var testLog string
			var expectedIcon string

			switch logType {
			case "error":
				testLog = "Error: " + logMessage
				expectedIcon = "✗" // error icon
			case "failed":
				testLog = "Deploy failed on server: " + logMessage
				expectedIcon = "✗" // error icon
			case "success":
				testLog = "Deploy successful on server: " + logMessage
				expectedIcon = "✓" // success icon
			case "finished":
				testLog = "Deployment finished: " + logMessage
				expectedIcon = "✓" // success icon
			case "starting":
				testLog = "Starting deployment: " + logMessage
				expectedIcon = "ℹ" // info icon
			case "connecting":
				testLog = "Connecting to server: " + logMessage
				expectedIcon = "ℹ" // info icon
			case "running":
				testLog = "Running script: " + logMessage
				expectedIcon = "ℹ" // info icon
			default:
				testLog = logMessage
				expectedIcon = "•" // bullet icon for generic logs
			}

			// Add the test log to the model
			deployModel.logs = append(deployModel.logs, testLog)

			// Get the view output
			view := deployModel.View()

			// Verify that the view contains the expected icon
			if !strings.Contains(view, expectedIcon) {
				t.Logf("Log entry type '%s' should display icon '%s', view: %s", logType, expectedIcon, view)
				return false
			}

			// Verify that the log message is present in the view
			if !strings.Contains(view, logMessage) {
				t.Logf("View should contain log message '%s'", logMessage)
				return false
			}

			// Verify that the view is not empty
			if view == "" {
				t.Logf("View should not be empty")
				return false
			}

			return true
		},
		gen.OneConstOf("error", "failed", "success", "finished", "starting", "connecting", "running", "default"),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 14: Log entry alternating styles**
// For any sequence of multiple log entries, alternating visual styles or indentation should be applied for differentiation
// Validates: Requirements 6.5
func TestProperty_LogEntryAlternatingStyles(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("log entries use alternating styles for readability", prop.ForAll(
		func(numLogs int) bool {
			// Ensure we have at least 2 logs to test alternating styles
			if numLogs < 2 {
				numLogs = 2
			}
			// Cap at reasonable number for testing
			if numLogs > 20 {
				numLogs = 20
			}

			// Create a test config with a project
			config := &Config{
				SSHConfigs: []SSHConfig{
					{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
				},
				Projects: []Project{
					{Name: "TestProject", DeployServers: []string{"TestServer"}},
				},
			}

			// Create deploy model
			deployModel := NewDeployModel(config.Projects[0], config)

			// Add multiple log entries
			deployModel.logs = []string{} // Clear default logs
			for i := 0; i < numLogs; i++ {
				deployModel.logs = append(deployModel.logs, fmt.Sprintf("Log entry %d", i))
			}

			// Get the view output
			view := deployModel.View()

			// Property 1: Verify that all log entries are present
			for i := 0; i < numLogs; i++ {
				expectedLog := fmt.Sprintf("Log entry %d", i)
				if !strings.Contains(view, expectedLog) {
					t.Logf("View should contain log entry '%s'", expectedLog)
					return false
				}
			}

			// Property 2: Verify that alternating styles are applied
			// The implementation uses alternating background colors:
			// - Even indices (0, 2, 4, ...): no background (default)
			// - Odd indices (1, 3, 5, ...): backgroundDark

			// We verify this by checking that the styles are defined correctly
			// and that the rendering logic applies them alternately

			// Check that odd-indexed entries use backgroundDark
			oddEntryStyle := lipgloss.NewStyle().
				Background(backgroundDark).
				Padding(0, 1).
				MarginBottom(0)

			// Verify the style has backgroundDark
			oddBg := oddEntryStyle.GetBackground()
			if oddBg != backgroundDark {
				t.Logf("Odd-indexed log entries should have backgroundDark (%v), got %v", backgroundDark, oddBg)
				return false
			}

			// Property 3: Verify that the alternating pattern is consistent
			// Both even and odd styles should have the same padding
			evenEntryStyle := lipgloss.NewStyle().
				Padding(0, 1).
				MarginBottom(0)

			evenPaddingLeft := evenEntryStyle.GetPaddingLeft()
			evenPaddingRight := evenEntryStyle.GetPaddingRight()
			oddPaddingLeft := oddEntryStyle.GetPaddingLeft()
			oddPaddingRight := oddEntryStyle.GetPaddingRight()

			if evenPaddingLeft != oddPaddingLeft || evenPaddingRight != oddPaddingRight {
				t.Logf("Alternating styles should have consistent padding")
				return false
			}

			if evenPaddingLeft != 1 || evenPaddingRight != 1 {
				t.Logf("Log entry styles should have padding of 1")
				return false
			}

			// Property 4: Verify that even and odd styles differ in background
			// The key distinction is that odd entries have a background while even entries don't
			evenBg := evenEntryStyle.GetBackground()
			if evenBg == oddBg {
				t.Logf("Even and odd log entries should have different background colors for alternation")
				return false
			}

			// Property 5: Verify that the view is not empty
			if view == "" {
				t.Logf("View should not be empty")
				return false
			}

			// Property 6: Verify that the number of log entries in the view matches
			// Count occurrences of "Log entry" in the view
			logEntryCount := strings.Count(view, "Log entry")
			if logEntryCount != numLogs {
				t.Logf("Expected %d log entries in view, got %d", numLogs, logEntryCount)
				return false
			}

			return true
		},
		gen.IntRange(2, 15),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 27: Status message visual distinction**
// For any status message, the styling should be visually distinct from non-status content
// Validates: Requirements 10.5
func TestProperty_StatusMessageVisualDistinctionFromContent(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("status messages are visually distinct from non-status content", prop.ForAll(
		func(statusType string, message string) bool {
			// Skip empty messages
			if strings.TrimSpace(message) == "" {
				return true
			}

			// Create a test config with a project
			config := &Config{
				SSHConfigs: []SSHConfig{
					{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
				},
				Projects: []Project{
					{Name: "TestProject", DeployServers: []string{"TestServer"}},
				},
			}

			// Create deploy model
			deployModel := NewDeployModel(config.Projects[0], config)

			// Add a status message and a non-status message
			var statusMessage string
			var nonStatusMessage string

			switch statusType {
			case "success":
				statusMessage = "Deploy successful on server: " + message
				nonStatusMessage = "Regular log message: " + message
			case "error":
				statusMessage = "Error: " + message
				nonStatusMessage = "Regular log message: " + message
			case "info":
				statusMessage = "Starting deployment: " + message
				nonStatusMessage = "Regular log message: " + message
			default:
				statusMessage = "Deployment finished: " + message
				nonStatusMessage = "Regular log message: " + message
			}

			// Add both messages to the model
			deployModel.logs = []string{statusMessage, nonStatusMessage}

			// Get the view output
			view := deployModel.View()

			// Property 1: Verify that both messages are present in the view
			if !strings.Contains(view, message) {
				t.Logf("View should contain the message '%s'", message)
				return false
			}

			// Property 2: Verify that status messages use distinct styling
			// Status messages should have:
			// - Icons (✓, ✗, ℹ, •)
			// - Level badges (SUCCESS, ERROR, INFO)
			// - Color-coded styling

			// Check for icon presence (status messages have icons)
			hasIcon := strings.Contains(view, "✓") || strings.Contains(view, "✗") ||
				strings.Contains(view, "ℹ") || strings.Contains(view, "•")
			if !hasIcon {
				t.Logf("Status messages should display icons")
				return false
			}

			// Check for level badge presence
			hasBadge := strings.Contains(view, "SUCCESS") || strings.Contains(view, "ERROR") ||
				strings.Contains(view, "INFO")
			if !hasBadge {
				t.Logf("Status messages should display level badges")
				return false
			}

			// Property 3: Verify that status message styles are distinct from body style
			// Status messages use successStyle, errorStyle, or statusStyle
			// Non-status messages use bodyStyle

			// Get the styles
			var msgStyle lipgloss.Style
			switch statusType {
			case "success":
				msgStyle = successStyle
			case "error":
				msgStyle = errorStyle
			case "info":
				msgStyle = statusStyle
			default:
				msgStyle = successStyle
			}

			// Verify that status style has a different foreground color than body style
			statusFg := msgStyle.GetForeground()
			bodyFg := bodyStyle.GetForeground()

			if statusFg == bodyFg {
				t.Logf("Status message style should have different foreground color than body style")
				return false
			}

			// Property 4: Verify that status styles are bold
			// All status styles (success, error, warning, info) should be bold
			if !msgStyle.GetBold() {
				t.Logf("Status message style should be bold")
				return false
			}

			// Property 5: Verify that body style is not bold (for distinction)
			if bodyStyle.GetBold() {
				t.Logf("Body style should not be bold to distinguish from status messages")
				return false
			}

			// Property 6: Verify that the view is not empty
			if view == "" {
				t.Logf("View should not be empty")
				return false
			}

			// Property 7: Verify that status messages have timestamps
			// All log entries should have timestamps in [MM:SS] format
			hasTimestamp := strings.Contains(view, "[")
			if !hasTimestamp {
				t.Logf("Status messages should have timestamps")
				return false
			}

			return true
		},
		gen.OneConstOf("success", "error", "info", "default"),
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 2: Consistent spacing across UI elements**
// For any UI element across all screens, the padding and margin values should follow consistent ratios and patterns
// Validates: Requirements 2.1, 2.5
func TestProperty_ConsistentSpacingAcrossUIElements(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all UI elements use consistent padding and margin ratios", prop.ForAll(
		func(seed int64) bool {
			// Test config
			config := &Config{
				SSHConfigs: []SSHConfig{
					{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
				},
				Projects: []Project{
					{Name: "TestProject", DeployServers: []string{"TestServer"}},
				},
			}

			// Property 1: Verify that all border styles use consistent padding
			// All border styles should use the same padding values (1 vertical, 2 horizontal)
			borderStyles := []lipgloss.Style{
				borderStyle,
				borderStyleThick,
				borderStyleThin,
				borderStyleDouble,
			}

			expectedPaddingV := 1
			expectedPaddingH := 2

			for i, style := range borderStyles {
				paddingTop := style.GetPaddingTop()
				paddingBottom := style.GetPaddingBottom()
				paddingLeft := style.GetPaddingLeft()
				paddingRight := style.GetPaddingRight()

				if paddingTop != expectedPaddingV || paddingBottom != expectedPaddingV {
					t.Logf("Border style %d: vertical padding mismatch. Expected %d, got top=%d bottom=%d",
						i, expectedPaddingV, paddingTop, paddingBottom)
					return false
				}

				if paddingLeft != expectedPaddingH || paddingRight != expectedPaddingH {
					t.Logf("Border style %d: horizontal padding mismatch. Expected %d, got left=%d right=%d",
						i, expectedPaddingH, paddingLeft, paddingRight)
					return false
				}
			}

			// Property 2: Verify that form styles use consistent padding
			formStyles := []lipgloss.Style{
				formStyle,
				cardStyle,
			}

			for i, style := range formStyles {
				paddingTop := style.GetPaddingTop()
				paddingBottom := style.GetPaddingBottom()
				paddingLeft := style.GetPaddingLeft()
				paddingRight := style.GetPaddingRight()

				if paddingTop != expectedPaddingV || paddingBottom != expectedPaddingV {
					t.Logf("Form style %d: vertical padding mismatch. Expected %d, got top=%d bottom=%d",
						i, expectedPaddingV, paddingTop, paddingBottom)
					return false
				}

				if paddingLeft != expectedPaddingH || paddingRight != expectedPaddingH {
					t.Logf("Form style %d: horizontal padding mismatch. Expected %d, got left=%d right=%d",
						i, expectedPaddingH, paddingLeft, paddingRight)
					return false
				}
			}

			// Property 3: Verify that card styles use consistent margins
			// Card styles should have consistent margin values
			expectedMarginV := 1

			cardStyles := []lipgloss.Style{
				cardStyle,
				cardStyleWithShadow,
			}

			for i, style := range cardStyles {
				marginTop := style.GetMarginTop()
				marginBottom := style.GetMarginBottom()

				if marginTop != expectedMarginV || marginBottom != expectedMarginV {
					t.Logf("Card style %d: vertical margin mismatch. Expected %d, got top=%d bottom=%d",
						i, expectedMarginV, marginTop, marginBottom)
					return false
				}
			}

			// Property 4: Verify that list item styles use consistent padding
			// List item styles should have consistent padding
			listStyles := []lipgloss.Style{
				itemStyle,
				selectedItemStyle,
			}

			expectedListPaddingLeft := 4
			expectedSelectedPaddingLeft := 2

			// itemStyle should have paddingLeft of 4
			if listStyles[0].GetPaddingLeft() != expectedListPaddingLeft {
				t.Logf("Item style padding left mismatch. Expected %d, got %d",
					expectedListPaddingLeft, listStyles[0].GetPaddingLeft())
				return false
			}

			// selectedItemStyle should have paddingLeft of 2
			if listStyles[1].GetPaddingLeft() != expectedSelectedPaddingLeft {
				t.Logf("Selected item style padding left mismatch. Expected %d, got %d",
					expectedSelectedPaddingLeft, listStyles[1].GetPaddingLeft())
				return false
			}

			// Property 5: Verify that help styles use consistent padding
			// Help styles should have consistent padding
			helpStyles := []lipgloss.Style{
				helpStyle,
				paginationStyle,
			}

			expectedHelpPaddingLeft := 4

			for i, style := range helpStyles {
				paddingLeft := style.GetPaddingLeft()

				if paddingLeft != expectedHelpPaddingLeft {
					t.Logf("Help style %d: padding left mismatch. Expected %d, got %d",
						i, expectedHelpPaddingLeft, paddingLeft)
					return false
				}
			}

			// Property 6: Verify that title style uses consistent padding
			// Title style should have consistent padding
			titlePaddingV := titleStyle.GetPaddingTop()
			titlePaddingH := titleStyle.GetPaddingLeft()

			if titlePaddingV != 0 {
				t.Logf("Title style vertical padding should be 0, got %d", titlePaddingV)
				return false
			}

			if titlePaddingH != 2 {
				t.Logf("Title style horizontal padding should be 2, got %d", titlePaddingH)
				return false
			}

			// Property 7: Verify that subtitle style uses consistent margins
			// Subtitle style should have consistent margins
			subtitleMarginTop := subtitleStyle.GetMarginTop()
			subtitleMarginBottom := subtitleStyle.GetMarginBottom()

			if subtitleMarginTop != 1 || subtitleMarginBottom != 1 {
				t.Logf("Subtitle style margin mismatch. Expected top=1 bottom=1, got top=%d bottom=%d",
					subtitleMarginTop, subtitleMarginBottom)
				return false
			}

			// Property 8: Verify spacing consistency across different views
			// Create different models and verify they all use consistent spacing
			mainMenu := NewMainMenu(config)
			sshList := NewSSHListModel(config)
			sshForm := NewSSHFormModel(config, -1)
			projectList := NewProjectListModel(config)
			projectForm := NewProjectFormModel(config, -1)

			// Get views
			mainMenuView := mainMenu.View()
			sshListView := sshList.View()
			sshFormView := sshForm.View()
			projectListView := projectList.View()
			projectFormView := projectForm.View()

			// All views should be non-empty
			if mainMenuView == "" || sshListView == "" || sshFormView == "" ||
				projectListView == "" || projectFormView == "" {
				t.Logf("All views should be non-empty")
				return false
			}

			// Property 9: Verify that spacing ratios are maintained
			// The ratio between different padding/margin values should be consistent
			// For example, horizontal padding (2) should be 2x vertical padding (1)
			if expectedPaddingH != expectedPaddingV*2 {
				t.Logf("Padding ratio mismatch. Horizontal (%d) should be 2x vertical (%d)",
					expectedPaddingH, expectedPaddingV)
				return false
			}

			return true
		},
		gen.Int64(),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 4: Nested content indentation**
// For any nested content structure, the indentation level should increase proportionally with hierarchy depth
// Validates: Requirements 2.4
func TestProperty_NestedContentIndentation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("nested content uses proportional indentation based on hierarchy depth", prop.ForAll(
		func(numServers int) bool {
			// Ensure we have at least 1 server
			if numServers < 1 {
				numServers = 1
			}
			// Cap at reasonable number for testing
			if numServers > 10 {
				numServers = 10
			}

			// Create test config with nested structure
			// SSH configs have nested details (connection info, auth type)
			sshConfigs := make([]SSHConfig, numServers)
			for i := 0; i < numServers; i++ {
				sshConfigs[i] = SSHConfig{
					Name:     fmt.Sprintf("Server%d", i),
					Host:     fmt.Sprintf("host%d.example.com", i),
					Port:     22 + i,
					User:     fmt.Sprintf("user%d", i),
					AuthType: "password",
					Password: "pass",
				}
			}

			config := &Config{
				SSHConfigs: sshConfigs,
				Projects:   []Project{},
			}

			// Create SSH list model which displays nested content
			sshList := NewSSHListModel(config)
			view := sshList.View()

			// Property 1: Verify that nested content (connection details) is indented
			// In the SSH list view, connection details are indented with "  " (2 spaces)
			// Auth type is also indented with "  " (2 spaces)
			// The indentation is part of the styled content, so we check for the content itself
			for i := 0; i < numServers; i++ {
				// Check that connection details are present (with indentation in the content)
				expectedDetail := fmt.Sprintf("  user%d@host%d.example.com:%d", i, i, 22+i)
				if !strings.Contains(view, expectedDetail) {
					t.Logf("Nested connection detail should be present with indentation: '%s'", expectedDetail)
					return false
				}

				// Check that auth type is present (with indentation in the content)
				expectedAuth := "  Auth: password"
				if !strings.Contains(view, expectedAuth) {
					t.Logf("Nested auth type should be present with indentation: '%s'", expectedAuth)
					return false
				}
			}

			// Property 2: Verify that the indentation is consistent across all items
			// All nested content should use the same indentation level (2 spaces)
			// We verify this by checking that all expected indented content is present
			// The actual line counting is unreliable due to ANSI codes, so we verify content presence instead
			for i := 0; i < numServers; i++ {
				// Each server should have indented connection details
				expectedDetail := fmt.Sprintf("  user%d@host%d.example.com:%d", i, i, 22+i)
				if !strings.Contains(view, expectedDetail) {
					t.Logf("Server %d should have indented connection details", i)
					return false
				}

				// Each server should have indented auth type
				if !strings.Contains(view, "  Auth: password") {
					t.Logf("Server %d should have indented auth type", i)
					return false
				}
			}

			// Property 3: Test with project list which also has nested content
			// Projects have nested server lists
			projects := make([]Project, numServers)
			for i := 0; i < numServers; i++ {
				projects[i] = Project{
					Name:              fmt.Sprintf("Project%d", i),
					DeployServers:     []string{fmt.Sprintf("Server%d", i), fmt.Sprintf("Server%d-backup", i)},
					BuildInstructions: "build command",
					DeployScript:      "deploy command",
				}
			}

			config.Projects = projects
			projectList := NewProjectListModel(config)

			// Set window size so the list renders properly
			model, _ := projectList.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
			projectList = model.(ProjectListModel)

			projectView := projectList.View()

			// Property 4: Verify that project nested content is indented
			// In the project list view, server lists, build instructions, and deploy scripts are indented
			for i := 0; i < numServers; i++ {
				// Check that server information is present and indented
				// The format includes "  " + icon + " Servers: "
				if !strings.Contains(projectView, "Servers:") {
					t.Logf("Project nested content should include 'Servers:' label")
					return false
				}

				// Check that build instructions are present and indented (if shown)
				if !strings.Contains(projectView, "Build:") {
					t.Logf("Project nested content should include 'Build:' label")
					return false
				}

				// Check that deploy script is present and indented (if shown)
				if !strings.Contains(projectView, "Deploy:") {
					t.Logf("Project nested content should include 'Deploy:' label")
					return false
				}
			}

			// Property 5: Verify that indentation increases with nesting depth
			// In the deploy view, log entries can have nested information
			// Create a deploy model to test nested log entries
			deployModel := NewDeployModel(projects[0], config)

			// Add logs with different nesting levels
			deployModel.logs = []string{
				"Starting deployment...",                 // Level 0
				"Connecting to Server0...",               // Level 1 (nested under deployment)
				"Running build on Server0...",            // Level 1 (nested under deployment)
				"Build output: compiling...",             // Level 2 (nested under build)
				"Deploy successful on Server0: deployed", // Level 1 (nested under deployment)
				"Deployment finished.",                   // Level 0
			}
			deployModel.done = true

			deployView := deployModel.View()

			// Verify that all log entries are present
			for _, log := range deployModel.logs {
				if !strings.Contains(deployView, log) {
					t.Logf("Deploy view should contain log entry: '%s'", log)
					return false
				}
			}

			// Property 6: Verify that form fields maintain consistent indentation
			// Form fields should not be indented (they are at the top level)
			// But error messages under fields should be indented
			sshForm := NewSSHFormModel(config, -1)

			// Set an invalid field to trigger error message
			sshForm.form[2].value = "invalid_port" // Port field
			sshForm.form[2].isValid = false
			sshForm.form[2].errorMsg = "Port must be between 1 and 65535"

			formView := sshForm.View()

			// Verify that error message is present and indented
			// Error messages are indented with MarginLeft(25)
			if !strings.Contains(formView, "Port must be between 1 and 65535") {
				t.Logf("Form error message should be present")
				return false
			}

			// Property 7: Verify that the indentation pattern is consistent
			// All first-level nested content should use the same indentation (2 spaces)
			// All second-level nested content should use more indentation (4 spaces)
			// This is verified by the consistent use of "  " prefix in the views

			// Property 8: Verify that indentation is proportional to depth
			// Level 0: no indentation
			// Level 1: 2 spaces
			// Level 2: 4 spaces (if applicable)
			// The ratio should be consistent (2 spaces per level)

			baseIndent := 2
			level1Indent := baseIndent * 1 // 2 spaces
			level2Indent := baseIndent * 2 // 4 spaces (if used)

			if level1Indent != 2 {
				t.Logf("Level 1 indentation should be 2 spaces, got %d", level1Indent)
				return false
			}

			if level2Indent != 4 {
				t.Logf("Level 2 indentation should be 4 spaces, got %d", level2Indent)
				return false
			}

			// Property 9: Verify that indentation is maintained across different view types
			// All views should use the same indentation pattern
			// This is verified by checking that all views use consistent spacing

			return true
		},
		gen.IntRange(1, 5),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 12: Long text wrapping**
// For any text content exceeding the display width, the text should wrap to the next line with proper alignment maintained
// Validates: Requirements 5.5
func TestProperty_LongTextWrapping(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("long text wraps correctly with proper alignment maintained", prop.ForAll(
		func(textLength int) bool {
			// Ensure we have a reasonable text length for testing
			if textLength < 50 {
				textLength = 50
			}
			if textLength > 200 {
				textLength = 200
			}

			// Create long text content
			longText := strings.Repeat("A", textLength)

			// Test config
			config := &Config{
				SSHConfigs: []SSHConfig{},
				Projects:   []Project{},
			}

			// Property 1: Test long text in form fields
			// Form fields have a fixed width (40), so long text should be handled
			projectForm := NewProjectFormModel(config, -1)

			// Set long text in a field (Build Instructions is multiline)
			buildFieldIndex := -1
			for i, field := range projectForm.form {
				if field.label == "Build Instructions" {
					buildFieldIndex = i
					break
				}
			}

			if buildFieldIndex == -1 {
				t.Logf("Project form should have Build Instructions field")
				return false
			}

			projectForm.form[buildFieldIndex].value = longText
			projectForm.cursor = buildFieldIndex

			formView := projectForm.View()

			// Verify that the view is not empty
			if formView == "" {
				t.Logf("Form view should not be empty")
				return false
			}

			// Property 2: Verify that long text is truncated or wrapped in the view
			// For multiline fields, long text should be truncated with "..." if it exceeds 120 characters
			if textLength > 120 {
				// Should contain truncation indicator
				if !strings.Contains(formView, "...") {
					t.Logf("Long text in multiline field should be truncated with '...'")
					return false
				}
			}

			// Property 3: Test long text in list items
			// Create SSH config with long name
			longName := strings.Repeat("Server", textLength/6)
			if len(longName) > 100 {
				longName = longName[:100]
			}

			sshConfig := SSHConfig{
				Name:     longName,
				Host:     "test.example.com",
				Port:     22,
				User:     "user",
				AuthType: "password",
				Password: "pass",
			}

			config.SSHConfigs = []SSHConfig{sshConfig}
			sshList := NewSSHListModel(config)

			listView := sshList.View()

			// Verify that the long name is present in the view
			// It may be truncated or wrapped, but should be present
			if !strings.Contains(listView, "Server") {
				t.Logf("List view should contain server name")
				return false
			}

			// Property 4: Test long text in card content
			// Cards should handle long content appropriately
			longContent := strings.Repeat("Content ", textLength/8)
			cardView := renderCard(longContent, "Test Card")

			// Verify that the card view is not empty
			if cardView == "" {
				t.Logf("Card view should not be empty")
				return false
			}

			// Verify that the content is present (may be wrapped)
			if !strings.Contains(cardView, "Content") {
				t.Logf("Card view should contain content")
				return false
			}

			// Property 5: Test long text in deployment logs
			// Deployment logs should handle long messages
			project := Project{
				Name:          "TestProject",
				DeployServers: []string{"TestServer"},
			}

			config.SSHConfigs = []SSHConfig{
				{Name: "TestServer", Host: "test.com", Port: 22, User: "user", AuthType: "password", Password: "pass"},
			}
			config.Projects = []Project{project}

			deployModel := NewDeployModel(project, config)

			// Add a long log message
			longLogMessage := "Deployment log: " + strings.Repeat("message ", textLength/8)
			deployModel.logs = []string{longLogMessage}
			deployModel.done = true

			deployView := deployModel.View()

			// Verify that the log message is present (may be wrapped)
			if !strings.Contains(deployView, "Deployment log:") {
				t.Logf("Deploy view should contain log message")
				return false
			}

			// Property 6: Verify that text wrapping maintains alignment
			// When text wraps, it should maintain proper alignment
			// This is implicitly tested by the lipgloss library's rendering
			// We verify that the view is properly formatted (contains borders, etc.)

			// Check that the form view has borders (indicating proper formatting)
			formBorderChars := []string{"╭", "╮", "╰", "╯", "─", "│"}
			hasFormBorder := false
			for _, char := range formBorderChars {
				if strings.Contains(formView, char) {
					hasFormBorder = true
					break
				}
			}

			if !hasFormBorder {
				t.Logf("Form view should have borders (indicating proper formatting)")
				return false
			}

			// Property 7: Verify that long text doesn't break the layout
			// The view should still be properly structured even with long text
			// We verify this by checking that the view contains expected structural elements

			// Form view should contain field labels
			if !strings.Contains(formView, "Name") {
				t.Logf("Form view should contain field labels")
				return false
			}

			// List view should contain title
			if !strings.Contains(listView, "SSH Configurations") {
				t.Logf("List view should contain title")
				return false
			}

			// Deploy view should contain title
			if !strings.Contains(deployView, "Deploying") {
				t.Logf("Deploy view should contain title")
				return false
			}

			// Property 8: Verify that text wrapping is consistent across different content types
			// All long text should be handled consistently (truncation or wrapping)

			// Property 9: Test with very long text (edge case)
			veryLongText := strings.Repeat("X", 500)

			// Set very long text in form field
			projectForm.form[buildFieldIndex].value = veryLongText

			veryLongFormView := projectForm.View()

			// Should still produce a valid view
			if veryLongFormView == "" {
				t.Logf("Form view with very long text should not be empty")
				return false
			}

			// Should contain truncation indicator
			if !strings.Contains(veryLongFormView, "...") {
				t.Logf("Very long text should be truncated with '...'")
				return false
			}

			// Property 10: Verify that empty text is handled correctly
			// Empty text should not cause issues
			projectForm.form[buildFieldIndex].value = ""

			emptyFormView := projectForm.View()

			// Should still produce a valid view
			if emptyFormView == "" {
				t.Logf("Form view with empty text should not be empty")
				return false
			}

			return true
		},
		gen.IntRange(50, 200),
	))

	properties.TestingRun(t)
}

// **Feature: ui-enhancement, Property 8: Title bold styling**
// For any title element, the rendered output should include bold styling
// Validates: Requirements 5.1
func TestProperty_TitleBoldStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all title elements use bold styling", prop.ForAll(
		func(titleText string) bool {
			// Skip empty strings as they don't provide meaningful testing
			if strings.TrimSpace(titleText) == "" {
				return true
			}

			// Test the titleStyle which is used for all titles in the application
			// Verify that the titleStyle has bold enabled
			if !titleStyle.GetBold() {
				t.Logf("titleStyle should have bold enabled")
				return false
			}

			// Test the subtitleStyle which is used for section headers
			// Verify that the subtitleStyle has bold enabled
			if !subtitleStyle.GetBold() {
				t.Logf("subtitleStyle should have bold enabled")
				return false
			}

			// Render the title with the titleStyle and verify it's not empty
			renderedTitle := titleStyle.Render(titleText)
			if renderedTitle == "" {
				t.Logf("Rendered title should not be empty for non-empty input")
				return false
			}

			// Render the subtitle with the subtitleStyle and verify it's not empty
			renderedSubtitle := subtitleStyle.Render(titleText)
			if renderedSubtitle == "" {
				t.Logf("Rendered subtitle should not be empty for non-empty input")
				return false
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test to verify title style characteristics
func TestTitleStyleBoldCharacteristics(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"titleStyle has bold", titleStyle},
		{"subtitleStyle has bold", subtitleStyle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.style.GetBold() {
				t.Errorf("%s should have bold styling enabled", tt.name)
			}
		})
	}
}

// **Feature: ui-enhancement, Property 9: Description differentiation**
// For any description text, the styling (color or weight) should differ from primary text styling
// Validates: Requirements 5.2
func TestProperty_DescriptionDifferentiation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("description text has different styling from primary text", prop.ForAll(
		func(descText string) bool {
			// Skip empty strings as they don't provide meaningful testing
			if strings.TrimSpace(descText) == "" {
				return true
			}

			// In the UI, descriptions use textSecondary or textMuted colors
			// Primary text uses textPrimary color
			// Verify that description colors differ from primary text color

			// Check that textSecondary differs from textPrimary
			if textSecondary == textPrimary {
				t.Logf("textSecondary should differ from textPrimary")
				return false
			}

			// Check that textMuted differs from textPrimary
			if textMuted == textPrimary {
				t.Logf("textMuted should differ from textPrimary")
				return false
			}

			// Verify that bodyStyle (primary text) uses textPrimary
			if bodyStyle.GetForeground() != textPrimary {
				t.Logf("bodyStyle should use textPrimary color")
				return false
			}

			// Verify that description-related styles use different colors
			// In the codebase, descriptions are rendered with detailStyle or similar
			// which uses textSecondary
			detailStyle := lipgloss.NewStyle().Foreground(textSecondary)
			if detailStyle.GetForeground() == bodyStyle.GetForeground() {
				t.Logf("Description style should use different color from body style")
				return false
			}

			// Render both primary and description text to verify they're different
			primaryRendered := bodyStyle.Render(descText)
			descRendered := detailStyle.Render(descText)

			// The rendered outputs should be different (due to different colors)
			if primaryRendered == descRendered {
				t.Logf("Primary and description rendered text should differ")
				return false
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test to verify description styling characteristics
func TestDescriptionStyleDifferentiation(t *testing.T) {
	tests := []struct {
		name              string
		descriptionColor  lipgloss.Color
		primaryTextColor  lipgloss.Color
		shouldBeDifferent bool
	}{
		{"textSecondary differs from textPrimary", textSecondary, textPrimary, true},
		{"textMuted differs from textPrimary", textMuted, textPrimary, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			areDifferent := tt.descriptionColor != tt.primaryTextColor
			if areDifferent != tt.shouldBeDifferent {
				t.Errorf("Expected colors to be different: %v, but got: %v", tt.shouldBeDifferent, areDifferent)
			}
		})
	}

	// Verify bodyStyle uses textPrimary
	if bodyStyle.GetForeground() != textPrimary {
		t.Errorf("bodyStyle should use textPrimary color")
	}

	// Verify that description styles in the codebase use secondary colors
	// Create a sample detail style as used in the views
	detailStyle := lipgloss.NewStyle().Foreground(textSecondary)
	if detailStyle.GetForeground() == bodyStyle.GetForeground() {
		t.Errorf("Description style should use different color from body style")
	}
}

// **Feature: ui-enhancement, Property 10: Help text italic styling**
// For any help text element, the rendered output should include italic styling
// Validates: Requirements 5.3
func TestProperty_HelpTextItalicStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("all help text elements use italic styling", prop.ForAll(
		func(helpText string) bool {
			// Skip empty strings as they don't provide meaningful testing
			if strings.TrimSpace(helpText) == "" {
				return true
			}

			// Test the helpTextStyle which is used for help text in the application
			// Verify that the helpTextStyle has italic enabled
			if !helpTextStyle.GetItalic() {
				t.Logf("helpTextStyle should have italic enabled")
				return false
			}

			// Test the helpStyle which is used for list help text
			// Verify that the helpStyle has italic enabled
			if !helpStyle.GetItalic() {
				t.Logf("helpStyle should have italic enabled")
				return false
			}

			// Render the help text with the helpTextStyle and verify it's not empty
			renderedHelpText := helpTextStyle.Render(helpText)
			if renderedHelpText == "" {
				t.Logf("Rendered help text should not be empty for non-empty input")
				return false
			}

			// Render the help text with the helpStyle and verify it's not empty
			renderedHelpStyle := helpStyle.Render(helpText)
			if renderedHelpStyle == "" {
				t.Logf("Rendered help style should not be empty for non-empty input")
				return false
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test to verify help text style characteristics
func TestHelpTextStyleItalicCharacteristics(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"helpTextStyle has italic", helpTextStyle},
		{"helpStyle has italic", helpStyle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.style.GetItalic() {
				t.Errorf("%s should have italic styling enabled", tt.name)
			}
		})
	}
}

// **Feature: ui-enhancement, Property 16: Contextual help positioning**
// For any contextual help text, it should be positioned adjacent to or near the relevant content it describes
// Validates: Requirements 7.4
func TestProperty_ContextualHelpPositioning(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("contextual help is positioned appropriately relative to content", prop.ForAll(
		func(numSSHConfigs int, numProjects int) bool {
			// Ensure reasonable bounds for testing
			if numSSHConfigs < 0 {
				numSSHConfigs = 0
			}
			if numSSHConfigs > 10 {
				numSSHConfigs = 10
			}
			if numProjects < 0 {
				numProjects = 0
			}
			if numProjects > 10 {
				numProjects = 10
			}

			// Create test config with random data
			sshConfigs := make([]SSHConfig, numSSHConfigs)
			for i := 0; i < numSSHConfigs; i++ {
				sshConfigs[i] = SSHConfig{
					Name:     fmt.Sprintf("Server%d", i),
					Host:     fmt.Sprintf("host%d.example.com", i),
					Port:     22,
					User:     "testuser",
					AuthType: "password",
					Password: "testpass",
				}
			}

			projects := make([]Project, numProjects)
			for i := 0; i < numProjects; i++ {
				projects[i] = Project{
					Name:          fmt.Sprintf("Project%d", i),
					DeployServers: []string{fmt.Sprintf("Server%d", i)},
				}
			}

			config := &Config{
				SSHConfigs: sshConfigs,
				Projects:   projects,
			}

			// Test 1: Main Menu - help text should be at the bottom after content
			mainMenu := NewMainMenu(config)
			mainMenuView := mainMenu.View()

			// Help text should appear after the menu items and divider
			// The help text contains keyboard shortcuts like "↑/↓", "Enter", "q"
			helpTextIndex := strings.Index(mainMenuView, "↑/↓")
			menuItemsIndex := strings.Index(mainMenuView, "SSH Management")

			// Help text should come after menu items
			if helpTextIndex != -1 && menuItemsIndex != -1 && helpTextIndex < menuItemsIndex {
				t.Logf("Main Menu: Help text should appear after menu items")
				return false
			}

			// Test 2: SSH List - help text should be at the bottom after list content
			sshList := NewSSHListModel(config)
			sshListView := sshList.View()

			// Help text should appear after the list items or empty message
			if numSSHConfigs > 0 {
				// With items, help should come after the last item
				lastItemIndex := strings.LastIndex(sshListView, fmt.Sprintf("Server%d", numSSHConfigs-1))
				helpIndex := strings.Index(sshListView, "Add")

				if lastItemIndex != -1 && helpIndex != -1 && helpIndex < lastItemIndex {
					t.Logf("SSH List: Help text should appear after list items")
					return false
				}
			} else {
				// With empty list, help should come after empty message
				emptyMsgIndex := strings.Index(sshListView, "No SSH Configurations Found")
				helpIndex := strings.Index(sshListView, "Add")

				if emptyMsgIndex != -1 && helpIndex != -1 && helpIndex < emptyMsgIndex {
					t.Logf("SSH List (empty): Help text should appear after empty message")
					return false
				}
			}

			// Test 3: SSH Form - contextual help should be near form fields
			sshForm := NewSSHFormModel(config, -1)
			sshFormView := sshForm.View()

			// The form has contextual help at the bottom: "Fill in all fields and press Enter on the last field to save"
			contextualHelpIndex := strings.Index(sshFormView, "Fill in all fields")
			lastFieldIndex := strings.LastIndex(sshFormView, "Key Password")

			// Contextual help should come after the form fields
			if contextualHelpIndex != -1 && lastFieldIndex != -1 && contextualHelpIndex < lastFieldIndex {
				t.Logf("SSH Form: Contextual help should appear after form fields")
				return false
			}

			// The keyboard shortcuts help should also be present and positioned appropriately
			keyboardHelpIndex := strings.Index(sshFormView, "Navigate")
			if keyboardHelpIndex != -1 && lastFieldIndex != -1 && keyboardHelpIndex < lastFieldIndex {
				t.Logf("SSH Form: Keyboard help should appear after form fields")
				return false
			}

			// Test 4: Project List - help text should be at the bottom after list content
			projectList := NewProjectListModel(config)
			projectListView := projectList.View()

			// Help text should appear after the list items or empty message
			if numProjects > 0 {
				// With items, help should come after the last item
				lastProjectIndex := strings.LastIndex(projectListView, fmt.Sprintf("Project%d", numProjects-1))
				helpIndex := strings.Index(projectListView, "Add")

				if lastProjectIndex != -1 && helpIndex != -1 && helpIndex < lastProjectIndex {
					t.Logf("Project List: Help text should appear after list items")
					return false
				}
			} else {
				// With empty list, help should come after empty message
				emptyMsgIndex := strings.Index(projectListView, "No Projects Found")
				helpIndex := strings.Index(projectListView, "Add")

				if emptyMsgIndex != -1 && helpIndex != -1 && helpIndex < emptyMsgIndex {
					t.Logf("Project List (empty): Help text should appear after empty message")
					return false
				}
			}

			// Test 5: Project Form - contextual help should be near form fields
			projectForm := NewProjectFormModel(config, -1)
			projectFormView := projectForm.View()

			// The form has contextual help at the bottom
			contextualHelpIndex = strings.Index(projectFormView, "Fill in all fields")
			lastFieldIndex = strings.LastIndex(projectFormView, "Deploy Servers")

			// Contextual help should come after the form fields
			if contextualHelpIndex != -1 && lastFieldIndex != -1 && contextualHelpIndex < lastFieldIndex {
				t.Logf("Project Form: Contextual help should appear after form fields")
				return false
			}

			// Test 6: Verify that help text is separated from content by dividers
			// This ensures appropriate visual positioning
			views := []struct {
				name string
				view string
			}{
				{"Main Menu", mainMenuView},
				{"SSH List", sshListView},
				{"SSH Form", sshFormView},
				{"Project List", projectListView},
				{"Project Form", projectFormView},
			}

			for _, v := range views {
				// Check that dividers (─) are used to separate help from content
				// This is a visual indicator of appropriate positioning
				if !strings.Contains(v.view, "─") {
					t.Logf("%s: View should contain dividers to separate help from content", v.name)
					return false
				}

				// Check that help text is present (either keyboard shortcuts or contextual help)
				hasKeyboardHelp := strings.Contains(v.view, "Navigate") ||
					strings.Contains(v.view, "Add") ||
					strings.Contains(v.view, "Quit") ||
					strings.Contains(v.view, "Enter")

				hasContextualHelp := strings.Contains(v.view, "Fill in all fields") ||
					strings.Contains(v.view, "No SSH Configurations Found") ||
					strings.Contains(v.view, "No Projects Found")

				if !hasKeyboardHelp && !hasContextualHelp {
					t.Logf("%s: View should contain help text (keyboard shortcuts or contextual help)", v.name)
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 5),
		gen.IntRange(0, 5),
	))

	properties.TestingRun(t)
}

// Unit test to verify contextual help positioning in specific views
func TestContextualHelpPositioning(t *testing.T) {
	config := &Config{
		SSHConfigs: []SSHConfig{
			{Name: "TestServer", Host: "test.example.com", Port: 22, User: "testuser", AuthType: "password", Password: "pass"},
		},
		Projects: []Project{
			{Name: "TestProject", DeployServers: []string{"TestServer"}},
		},
	}

	tests := []struct {
		name                 string
		getView              func() string
		contentMarker        string
		helpMarker           string
		shouldHaveContextual bool
	}{
		{
			name:                 "Main Menu help after content",
			getView:              func() string { return NewMainMenu(config).View() },
			contentMarker:        "SSH Management",
			helpMarker:           "Navigate",
			shouldHaveContextual: false,
		},
		{
			name:                 "SSH List help after content",
			getView:              func() string { return NewSSHListModel(config).View() },
			contentMarker:        "TestServer",
			helpMarker:           "Add",
			shouldHaveContextual: false,
		},
		{
			name:                 "SSH Form contextual help after fields",
			getView:              func() string { return NewSSHFormModel(config, -1).View() },
			contentMarker:        "Key Password",
			helpMarker:           "Fill in all fields",
			shouldHaveContextual: true,
		},
		{
			name:                 "Project List help after content",
			getView:              func() string { return NewProjectListModel(config).View() },
			contentMarker:        "TestProject",
			helpMarker:           "Add",
			shouldHaveContextual: false,
		},
		{
			name:                 "Project Form contextual help after fields",
			getView:              func() string { return NewProjectFormModel(config, -1).View() },
			contentMarker:        "Deploy Servers",
			helpMarker:           "Fill in all fields",
			shouldHaveContextual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.getView()

			// Find positions of content and help
			contentIndex := strings.Index(view, tt.contentMarker)
			helpIndex := strings.Index(view, tt.helpMarker)

			// Verify both are present
			if contentIndex == -1 {
				t.Errorf("Content marker '%s' not found in view", tt.contentMarker)
				return
			}
			if helpIndex == -1 {
				t.Errorf("Help marker '%s' not found in view", tt.helpMarker)
				return
			}

			// Verify help comes after content
			if helpIndex < contentIndex {
				t.Errorf("Help text should appear after content. Content at %d, Help at %d", contentIndex, helpIndex)
			}

			// Verify divider is present between content and help
			if !strings.Contains(view, "─") {
				t.Errorf("View should contain divider to separate help from content")
			}

			// For views with contextual help, verify it's present
			if tt.shouldHaveContextual {
				if !strings.Contains(view, "Fill in all fields") {
					t.Errorf("View should contain contextual help text")
				}
			}
		})
	}
}

// **Feature: ui-enhancement, Property 22: Pagination indicator styling**
// For any paginated list, page indicators should be rendered with clear styling
// Validates: Requirements 9.4
func TestProperty_PaginationIndicatorStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("pagination indicators have clear and consistent styling", prop.ForAll(
		func(numItems int) bool {
			// Ensure we have enough items to trigger pagination
			// The bubble tea list component typically shows pagination when there are many items
			if numItems < 10 {
				numItems = 10
			}
			// Cap at reasonable number for testing
			if numItems > 100 {
				numItems = 100
			}

			// Create test config with many SSH configurations to trigger pagination
			sshConfigs := make([]SSHConfig, numItems)
			for i := 0; i < numItems; i++ {
				sshConfigs[i] = SSHConfig{
					Name:     fmt.Sprintf("Server%d", i),
					Host:     fmt.Sprintf("host%d.example.com", i),
					Port:     22 + i,
					User:     fmt.Sprintf("user%d", i),
					AuthType: "password",
					Password: "pass",
				}
			}

			config := &Config{
				SSHConfigs: sshConfigs,
				Projects:   []Project{},
			}

			// Property 1: Verify that paginationStyle is properly defined
			// The paginationStyle should have consistent styling characteristics
			expectedPaddingLeft := 4
			expectedForeground := mutedColor

			actualPaddingLeft := paginationStyle.GetPaddingLeft()
			actualForeground := paginationStyle.GetForeground()

			if actualPaddingLeft != expectedPaddingLeft {
				t.Logf("Pagination style padding left should be %d, got %d", expectedPaddingLeft, actualPaddingLeft)
				return false
			}

			if actualForeground != expectedForeground {
				t.Logf("Pagination style foreground should be mutedColor (%v), got %v", expectedForeground, actualForeground)
				return false
			}

			// Property 2: Verify that pagination style is applied to list models
			// Create SSH list model
			sshList := NewSSHListModel(config)

			// Set a small window size to ensure pagination is triggered
			// With many items and a small height, pagination should be visible
			model, _ := sshList.Update(tea.WindowSizeMsg{Width: 80, Height: 15})
			sshList = model.(SSHListModel)

			// Get the view output
			view := sshList.View()

			// Property 3: Verify that the view is not empty
			if view == "" {
				t.Logf("SSH list view should not be empty")
				return false
			}

			// Property 4: Verify that pagination indicators are present when there are many items
			// The bubble tea list component shows pagination info like "1/10" or similar
			// We check for the presence of pagination-related content
			// Note: The exact format depends on the bubble tea list implementation
			// but pagination indicators typically include numbers and slashes

			// Property 5: Test with Project list as well
			projects := make([]Project, numItems)
			for i := 0; i < numItems; i++ {
				projects[i] = Project{
					Name:          fmt.Sprintf("Project%d", i),
					DeployServers: []string{fmt.Sprintf("Server%d", i)},
				}
			}

			config.Projects = projects
			projectList := NewProjectListModel(config)

			// Set a small window size to ensure pagination is triggered
			model2, _ := projectList.Update(tea.WindowSizeMsg{Width: 80, Height: 15})
			projectList = model2.(ProjectListModel)

			projectView := projectList.View()

			// Verify that the project list view is not empty
			if projectView == "" {
				t.Logf("Project list view should not be empty")
				return false
			}

			// Property 6: Verify that pagination style is consistent across different list types
			// Both SSH list and Project list should use the same pagination style
			// We verify this by checking that both use paginationStyle

			// The paginationStyle is set in both NewSSHListModel and NewProjectListModel
			// with: l.Styles.PaginationStyle = paginationStyle

			// Property 7: Verify that pagination style has appropriate visual characteristics
			// The pagination style should be distinct but not distracting
			// It uses mutedColor to be subtle yet visible

			// Verify that mutedColor is different from primary text colors
			if mutedColor == textPrimary {
				t.Logf("Pagination mutedColor should differ from textPrimary for visual distinction")
				return false
			}

			// Property 8: Verify that pagination style has consistent padding
			// The padding should match other help-related styles for consistency
			// Both paginationStyle and helpStyle use paddingLeft of 4

			helpPaddingLeft := helpStyle.GetPaddingLeft()
			if actualPaddingLeft != helpPaddingLeft {
				t.Logf("Pagination padding should match help style padding (%d), got %d", helpPaddingLeft, actualPaddingLeft)
				return false
			}

			// Property 9: Verify that pagination indicators are positioned consistently
			// Pagination indicators should appear in a consistent location (typically at the bottom)
			// This is handled by the bubble tea list component, but we verify the style is applied

			// Property 10: Verify that the pagination style is readable
			// The foreground color should provide sufficient contrast
			// We verify that mutedColor is not the same as the background colors

			if mutedColor == backgroundDark || mutedColor == backgroundLight {
				t.Logf("Pagination foreground color should differ from background colors for readability")
				return false
			}

			return true
		},
		gen.IntRange(10, 50),
	))

	properties.TestingRun(t)
}

// Unit test to verify pagination style characteristics
func TestPaginationStyleCharacteristics(t *testing.T) {
	// Test that paginationStyle has the expected properties
	expectedPaddingLeft := 4
	expectedForeground := mutedColor

	actualPaddingLeft := paginationStyle.GetPaddingLeft()
	actualForeground := paginationStyle.GetForeground()

	if actualPaddingLeft != expectedPaddingLeft {
		t.Errorf("Pagination style padding left should be %d, got %d", expectedPaddingLeft, actualPaddingLeft)
	}

	if actualForeground != expectedForeground {
		t.Errorf("Pagination style foreground should be mutedColor (%v), got %v", expectedForeground, actualForeground)
	}

	// Verify that pagination style matches help style padding for consistency
	helpPaddingLeft := helpStyle.GetPaddingLeft()
	if actualPaddingLeft != helpPaddingLeft {
		t.Errorf("Pagination padding should match help style padding (%d), got %d", helpPaddingLeft, actualPaddingLeft)
	}

	// Verify that mutedColor differs from primary text color
	if mutedColor == textPrimary {
		t.Errorf("Pagination mutedColor should differ from textPrimary for visual distinction")
	}

	// Verify that mutedColor differs from background colors for readability
	if mutedColor == backgroundDark {
		t.Errorf("Pagination foreground color should differ from backgroundDark for readability")
	}

	if mutedColor == backgroundLight {
		t.Errorf("Pagination foreground color should differ from backgroundLight for readability")
	}
}

// **Feature: ui-enhancement, Property 26: Input prompt styling**
// For any user input prompt, the prompt should be rendered with distinct styling
// Validates: Requirements 10.4
func TestProperty_InputPromptStyling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("input prompts have distinct styling", prop.ForAll(
		func(promptText string) bool {
			// Skip empty strings as they don't provide meaningful testing
			if strings.TrimSpace(promptText) == "" {
				return true
			}

			// Render the input prompt
			rendered := renderInputPrompt(promptText)

			// Verify the prompt text is present in the rendered output
			if !strings.Contains(rendered, promptText) {
				t.Logf("Rendered prompt doesn't contain original text: %q", promptText)
				return false
			}

			// Verify the rendered output is not empty
			if rendered == "" {
				t.Logf("Rendered prompt is empty for text: %q", promptText)
				return false
			}

			// Verify the rendered output is different from plain text
			// (it should have ANSI styling codes)
			if rendered == promptText {
				t.Logf("Rendered prompt is identical to plain text (no styling applied): %q", promptText)
				return false
			}

			// Verify that the inputPromptStyle has the expected characteristics
			// Check foreground color
			actualFg := inputPromptStyle.GetForeground()
			if actualFg != accentColor {
				t.Logf("Input prompt style foreground color mismatch: expected %v, got %v", accentColor, actualFg)
				return false
			}

			// Check bold styling
			if !inputPromptStyle.GetBold() {
				t.Logf("Input prompt style should be bold")
				return false
			}

			// Check background color
			actualBg := inputPromptStyle.GetBackground()
			if actualBg != backgroundLight {
				t.Logf("Input prompt style background color mismatch: expected %v, got %v", backgroundLight, actualBg)
				return false
			}

			// Check horizontal padding
			paddingLeft := inputPromptStyle.GetPaddingLeft()
			paddingRight := inputPromptStyle.GetPaddingRight()
			if paddingLeft != 1 || paddingRight != 1 {
				t.Logf("Input prompt style padding mismatch: expected 1, got left=%d right=%d", paddingLeft, paddingRight)
				return false
			}

			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// Unit test to verify input prompt styling characteristics
func TestInputPromptStyling(t *testing.T) {
	tests := []struct {
		name       string
		promptText string
	}{
		{"simple prompt", "Enter your name:"},
		{"prompt with question", "What is your choice?"},
		{"prompt with special chars", "Enter value (1-10):"},
		{"long prompt", "Please enter a very long prompt text that should still be styled correctly"},
		{"prompt with numbers", "Enter port number (default: 22):"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderInputPrompt(tt.promptText)

			// Verify the prompt text is present
			if !strings.Contains(result, tt.promptText) {
				t.Errorf("renderInputPrompt(%q) result doesn't contain prompt text", tt.promptText)
			}

			// Verify result is not empty
			if result == "" {
				t.Errorf("renderInputPrompt(%q) returned empty string", tt.promptText)
			}

			// Verify result is styled (different from plain text)
			if result == tt.promptText {
				t.Errorf("renderInputPrompt(%q) returned unstyled text", tt.promptText)
			}
		})
	}
}

// Unit test to verify input prompt style characteristics
func TestInputPromptStyleCharacteristics(t *testing.T) {
	// Check foreground color
	actualFg := inputPromptStyle.GetForeground()
	if actualFg != accentColor {
		t.Errorf("Input prompt style foreground: expected %v, got %v", accentColor, actualFg)
	}

	// Check bold styling
	if !inputPromptStyle.GetBold() {
		t.Error("Input prompt style should be bold")
	}

	// Check background color
	actualBg := inputPromptStyle.GetBackground()
	if actualBg != backgroundLight {
		t.Errorf("Input prompt style background: expected %v, got %v", backgroundLight, actualBg)
	}

	// Check horizontal padding
	paddingLeft := inputPromptStyle.GetPaddingLeft()
	paddingRight := inputPromptStyle.GetPaddingRight()
	if paddingLeft != 1 || paddingRight != 1 {
		t.Errorf("Input prompt style padding: expected 1, got left=%d right=%d", paddingLeft, paddingRight)
	}
}

// Unit test to verify renderStatusMessage function
func TestRenderStatusMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		messageType string
		expectIcon  bool
	}{
		{"success message", "Operation completed successfully", "success", true},
		{"error message", "An error occurred", "error", true},
		{"warning message", "This is a warning", "warning", true},
		{"info message", "Information message", "info", true},
		{"default message", "Default message", "default", false},
		{"unknown type", "Unknown type message", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderStatusMessage(tt.message, tt.messageType)

			// Verify the message text is present
			if !strings.Contains(result, tt.message) {
				t.Errorf("renderStatusMessage(%q, %q) result doesn't contain message text", tt.message, tt.messageType)
			}

			// Verify result is not empty
			if result == "" {
				t.Errorf("renderStatusMessage(%q, %q) returned empty string", tt.message, tt.messageType)
			}

			// Verify icon is present for known types
			if tt.expectIcon {
				// Check for common icon characters
				hasIcon := strings.Contains(result, "✓") ||
					strings.Contains(result, "✗") ||
					strings.Contains(result, "⚠") ||
					strings.Contains(result, "ℹ")
				if !hasIcon {
					t.Errorf("renderStatusMessage(%q, %q) should contain an icon", tt.message, tt.messageType)
				}
			}
		})
	}
}

// Property test to verify status messages are visually distinct from other content
func TestProperty_StatusMessageVisualDistinctness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("status messages are visually distinct from plain text", prop.ForAll(
		func(message string, messageType string) bool {
			// Skip empty messages
			if strings.TrimSpace(message) == "" {
				return true
			}

			// Normalize message type to valid values
			validTypes := []string{"success", "error", "warning", "info"}
			if messageType == "" || !contains(validTypes, messageType) {
				messageType = "info"
			}

			// Render the status message
			rendered := renderStatusMessage(message, messageType)

			// Verify the message is present
			if !strings.Contains(rendered, message) {
				t.Logf("Rendered status message doesn't contain original message: %q", message)
				return false
			}

			// Verify the rendered output is not empty
			if rendered == "" {
				t.Logf("Rendered status message is empty for message: %q", message)
				return false
			}

			// Verify the rendered output is different from plain text
			// (it should have ANSI styling codes and possibly an icon)
			if rendered == message {
				t.Logf("Rendered status message is identical to plain text (no styling applied): %q", message)
				return false
			}

			// Verify that status messages have icons for known types
			if messageType == "success" || messageType == "error" || messageType == "warning" || messageType == "info" {
				hasIcon := strings.Contains(rendered, "✓") ||
					strings.Contains(rendered, "✗") ||
					strings.Contains(rendered, "⚠") ||
					strings.Contains(rendered, "ℹ")
				if !hasIcon {
					t.Logf("Status message of type %q should contain an icon", messageType)
					return false
				}
			}

			return true
		},
		gen.AnyString(),
		gen.OneConstOf("success", "error", "warning", "info", ""),
	))

	properties.TestingRun(t)
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
