package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/diiyw/ed/ssh"
)

// Enhanced Color Palette
var (
	// Primary colors for main UI elements
	primaryColor    = lipgloss.Color("#7C3AED") // Vibrant purple
	secondaryColor  = lipgloss.Color("#10B981") // Emerald green
	accentColor     = lipgloss.Color("#F59E0B") // Amber
	backgroundDark  = lipgloss.Color("#1F2937") // Dark gray
	backgroundLight = lipgloss.Color("#374151") // Medium gray

	// Status colors
	successColor = lipgloss.Color("#10B981") // Green
	errorColor   = lipgloss.Color("#EF4444") // Red
	warningColor = lipgloss.Color("#F59E0B") // Amber
	infoColor    = lipgloss.Color("#3B82F6") // Blue
	mutedColor   = lipgloss.Color("#9CA3AF") // Gray

	// Text colors
	textPrimary   = lipgloss.Color("#F9FAFB") // Almost white
	textSecondary = lipgloss.Color("#D1D5DB") // Light gray
	textMuted     = lipgloss.Color("#9CA3AF") // Medium gray
)

// Base Typography Styles
var (
	// Title style: Large, bold, centered for main titles
	titleStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(primaryColor).
			Padding(0, 2).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Align(lipgloss.Center)

	// Subtitle style: Medium, semi-bold for section headers
	subtitleStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	// Body style: Regular weight for main content
	bodyStyle = lipgloss.NewStyle().
			Foreground(textPrimary)

	// Label style: Small, uppercase for form labels
	labelStyle = lipgloss.NewStyle().
			Foreground(textSecondary).
			Bold(true)

	// Help style: Italic, muted for help text
	helpTextStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true)

	// Monospace style: For technical details (host, port, etc.)
	monospaceStyle = lipgloss.NewStyle().
			Foreground(textSecondary)
)

// Component Styles (updated to use new color palette)
var (
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(textPrimary)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(accentColor).
				Bold(true).
				Background(backgroundLight)

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4).
			Foreground(mutedColor)

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(textMuted).
			Italic(true)

	// Enhanced Border Style Variants
	// Default border style with rounded corners
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Thick border variant for emphasis
	borderStyleThick = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2)

	// Thin border variant for subtle separation
	borderStyleThin = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2)

	// Double border variant for special emphasis
	borderStyleDouble = lipgloss.NewStyle().
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(accentColor).
				Padding(1, 2)

	// Enhanced form style with rounded border and shadow effect
	formStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2).
			Width(80)

	// Card style with shadow effect using box characters
	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Card style with shadow effect (alternative)
	cardStyleWithShadow = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				BorderBottom(true).
				BorderRight(true).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(1).
				MarginRight(1)

	// Divider styles with decorative elements
	dividerStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	dividerStyleDecorative = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	dividerStyleAccent = lipgloss.NewStyle().
				Foreground(accentColor)

	statusStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// Input prompt style: Distinct styling for user input prompts
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Background(backgroundLight).
				Padding(0, 1)
)

// UI Helper Functions

// renderIcon returns a styled icon/symbol based on the icon type
func renderIcon(iconType string, style lipgloss.Style) string {
	icons := map[string]string{
		"success":  "✓",
		"error":    "✗",
		"warning":  "⚠",
		"info":     "ℹ",
		"loading":  "⋯",
		"arrow":    "→",
		"bullet":   "•",
		"check":    "✓",
		"cross":    "✗",
		"star":     "★",
		"heart":    "♥",
		"folder":   "📁",
		"file":     "📄",
		"server":   "🖥",
		"key":      "🔑",
		"lock":     "🔒",
		"unlock":   "🔓",
		"user":     "👤",
		"clock":    "🕐",
		"rocket":   "🚀",
		"gear":     "⚙",
		"wrench":   "🔧",
		"package":  "📦",
		"download": "⬇",
		"upload":   "⬆",
		"play":     "▶",
		"pause":    "⏸",
		"stop":     "⏹",
		"refresh":  "↻",
		"home":     "🏠",
		"search":   "🔍",
		"edit":     "✎",
		"delete":   "🗑",
		"add":      "➕",
		"remove":   "➖",
		"circle":   "●",
		"dot":      "·",
		"diamond":  "◆",
		"square":   "■",
		"triangle": "▲",
	}

	icon, exists := icons[iconType]
	if !exists {
		icon = iconType // Use the provided string if not found
	}

	return style.Render(icon)
}

// renderBadge creates a pill-shaped badge with appropriate styling
func renderBadge(text string, badgeType string) string {
	var style lipgloss.Style

	switch badgeType {
	case "success":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(successColor).
			Padding(0, 1).
			Bold(true)
	case "error":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(errorColor).
			Padding(0, 1).
			Bold(true)
	case "warning":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(warningColor).
			Padding(0, 1).
			Bold(true)
	case "info":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(infoColor).
			Padding(0, 1).
			Bold(true)
	case "primary":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true)
	case "secondary":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(secondaryColor).
			Padding(0, 1).
			Bold(true)
	case "muted":
		style = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(mutedColor).
			Padding(0, 1)
	default:
		style = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(backgroundLight).
			Padding(0, 1)
	}

	return style.Render(text)
}

// renderDivider creates a decorative horizontal line
func renderDivider(width int, style lipgloss.Style) string {
	if width <= 0 {
		width = 40
	}
	divider := strings.Repeat("─", width)
	return style.Render(divider)
}

// renderDividerDecorative creates a decorative horizontal line with ornamental elements
func renderDividerDecorative(width int, style lipgloss.Style, decorationType string) string {
	if width <= 0 {
		width = 40
	}

	var divider string
	switch decorationType {
	case "dots":
		// Divider with dots: ·─────·
		if width >= 4 {
			divider = "·" + strings.Repeat("─", width-2) + "·"
		} else {
			divider = strings.Repeat("─", width)
		}
	case "arrows":
		// Divider with arrows: ◄─────►
		if width >= 4 {
			divider = "◄" + strings.Repeat("─", width-2) + "►"
		} else {
			divider = strings.Repeat("─", width)
		}
	case "stars":
		// Divider with stars: ★─────★
		if width >= 4 {
			divider = "★" + strings.Repeat("─", width-2) + "★"
		} else {
			divider = strings.Repeat("─", width)
		}
	case "diamonds":
		// Divider with diamonds: ◆─────◆
		if width >= 4 {
			divider = "◆" + strings.Repeat("─", width-2) + "◆"
		} else {
			divider = strings.Repeat("─", width)
		}
	case "double":
		// Double line divider: ═════
		divider = strings.Repeat("═", width)
	case "wave":
		// Wave pattern: ～～～～～
		divider = strings.Repeat("～", width)
	default:
		// Default simple divider
		divider = strings.Repeat("─", width)
	}

	return style.Render(divider)
}

// renderProgressBar creates a progress indicator
func renderProgressBar(current, total int, width int) string {
	if total <= 0 {
		total = 1
	}
	if current > total {
		current = total
	}
	if width <= 0 {
		width = 40
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	if filled > width {
		filled = width
	}

	percentText := fmt.Sprintf(" %d%%", int(percentage*100))

	barStyle := lipgloss.NewStyle().Foreground(secondaryColor)
	emptyStyle := lipgloss.NewStyle().Foreground(mutedColor)

	styledBar := barStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", width-filled))

	return styledBar + percentText
}

// renderSpinner returns animated spinner frame
func renderSpinner(frame int) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	index := frame % len(spinners)
	return statusStyle.Render(spinners[index])
}

// renderCard wraps content in a styled card with optional title
func renderCard(content string, title string) string {
	if title != "" {
		titleBar := subtitleStyle.Render(title)
		return cardStyle.Render(titleBar + "\n\n" + content)
	}

	return cardStyle.Render(content)
}

// renderCardWithShadow wraps content in a styled card with shadow effect
func renderCardWithShadow(content string, title string) string {
	var cardContent string
	if title != "" {
		titleBar := subtitleStyle.Render(title)
		cardContent = cardStyleWithShadow.Render(titleBar + "\n\n" + content)
	} else {
		cardContent = cardStyleWithShadow.Render(content)
	}

	// Add shadow effect using box drawing characters
	lines := strings.Split(cardContent, "\n")
	var result strings.Builder
	for i, line := range lines {
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	// Add bottom shadow line
	if len(lines) > 0 {
		// Calculate approximate width of the card
		maxWidth := 0
		for _, line := range lines {
			// Strip ANSI codes for accurate width calculation (simplified)
			if len(line) > maxWidth {
				maxWidth = len(line)
			}
		}
		if maxWidth > 2 {
			shadowLine := " " + strings.Repeat("▀", maxWidth-2)
			result.WriteString("\n")
			result.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render(shadowLine))
		}
	}

	return result.String()
}

// renderKeyHelp formats keyboard shortcuts beautifully
func renderKeyHelp(keys map[string]string) string {
	if len(keys) == 0 {
		return ""
	}

	var parts []string
	keyStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(textSecondary)

	for key, desc := range keys {
		part := keyStyle.Render(key) + descStyle.Render(": "+desc)
		parts = append(parts, part)
	}

	separator := descStyle.Render(" • ")
	return strings.Join(parts, separator)
}

// renderInputPrompt formats an input prompt with distinct styling
func renderInputPrompt(promptText string) string {
	return inputPromptStyle.Render(promptText)
}

// renderStatusMessage formats a status message with appropriate styling based on message type
func renderStatusMessage(message string, messageType string) string {
	var style lipgloss.Style
	var icon string

	switch messageType {
	case "success":
		style = successStyle
		icon = renderIcon("success", successStyle)
	case "error":
		style = errorStyle
		icon = renderIcon("error", errorStyle)
	case "warning":
		style = warningStyle
		icon = renderIcon("warning", warningStyle)
	case "info":
		style = statusStyle
		icon = renderIcon("info", statusStyle)
	default:
		style = bodyStyle
		icon = ""
	}

	if icon != "" {
		return icon + " " + style.Render(message)
	}
	return style.Render(message)
}

func getCustomDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = itemStyle
	d.Styles.NormalDesc = itemStyle
	d.Styles.SelectedTitle = selectedItemStyle
	d.Styles.SelectedDesc = selectedItemStyle
	return d
}

// Main Menu Types
type mainMenuItem struct {
	title, desc string
}

func (i mainMenuItem) Title() string       { return i.title }
func (i mainMenuItem) Description() string { return i.desc }
func (i mainMenuItem) FilterValue() string { return i.title }

// MainMenuModel represents the main menu
type MainMenuModel struct {
	list     list.Model
	choice   string
	quitting bool
	config   *Config
}

// NewMainMenu creates a new main menu
func NewMainMenu(config *Config) MainMenuModel {
	items := []list.Item{
		mainMenuItem{title: "SSH Management", desc: "Add, edit, delete, and test SSH server configurations"},
		mainMenuItem{title: "Project Management", desc: "Manage projects and deployments"},
		mainMenuItem{title: "Exit", desc: "Quit the application"},
	}

	l := list.New(items, getCustomDelegate(), 0, 0)
	l.Title = "Easy Deploy - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return MainMenuModel{
		list:   l,
		config: config,
	}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := titleStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(mainMenuItem)
			if ok {
				m.choice = i.title
				switch i.title {
				case "SSH Management":
					return NewSSHListModel(m.config), nil
				case "Project Management":
					return NewProjectListModel(m.config), nil
				case "Exit":
					m.quitting = true
					return m, tea.Quit
				}
			}
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) View() string {
	if m.quitting {
		return successStyle.Render("Goodbye!\n")
	}

	var content strings.Builder

	// Welcome banner with decorative header
	bannerStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(primaryColor).
		Padding(1, 4).
		Bold(true).
		Align(lipgloss.Center).
		Width(60)

	welcomeBanner := bannerStyle.Render("🚀 EASY DEPLOY 🚀")
	content.WriteString(welcomeBanner + "\n\n")

	// Subtitle with decorative divider
	subtitleText := subtitleStyle.Render("Main Menu")
	content.WriteString(subtitleText + "\n")
	content.WriteString(renderDividerDecorative(60, dividerStyleAccent, "stars") + "\n\n")

	// Enhanced menu items with icons and improved selection indicators
	menuItems := []struct {
		title string
		desc  string
		icon  string
	}{
		{"SSH Management", "Add, edit, delete, and test SSH server configurations", "server"},
		{"Project Management", "Manage projects and deployments", "package"},
		{"Exit", "Quit the application", "cross"},
	}

	selectedIndex := m.list.Index()

	for i, item := range menuItems {
		// Icon for menu item
		iconStyle := lipgloss.NewStyle().Foreground(primaryColor)
		itemIcon := renderIcon(item.icon, iconStyle)

		// Selection indicator
		var indicator string
		var itemTitleStyle lipgloss.Style
		var itemDescStyle lipgloss.Style
		var itemContainer lipgloss.Style

		if i == selectedIndex {
			// Selected item with highlight box
			indicator = renderIcon("arrow", lipgloss.NewStyle().Foreground(accentColor).Bold(true))
			itemTitleStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				Bold(true)
			itemDescStyle = lipgloss.NewStyle().
				Foreground(textSecondary)
			itemContainer = lipgloss.NewStyle().
				Background(backgroundLight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(1, 2).
				Width(56).
				MarginBottom(1)
		} else {
			// Unselected item
			indicator = "  "
			itemTitleStyle = lipgloss.NewStyle().
				Foreground(textPrimary)
			itemDescStyle = lipgloss.NewStyle().
				Foreground(textMuted)
			itemContainer = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Padding(1, 2).
				Width(56).
				MarginBottom(1)
		}

		// Build menu item content
		itemContent := indicator + " " + itemIcon + " " + itemTitleStyle.Render(item.title) + "\n" +
			"    " + itemDescStyle.Render(item.desc)

		content.WriteString(itemContainer.Render(itemContent) + "\n")

		// Add visual separator between menu sections (but not after the last item)
		if i < len(menuItems)-1 {
			separatorStyle := lipgloss.NewStyle().Foreground(mutedColor)
			separator := separatorStyle.Render(strings.Repeat("·", 60))
			content.WriteString(separator + "\n\n")
		}
	}

	// Add decorative divider before help text
	content.WriteString("\n" + renderDividerDecorative(60, dividerStyle, "default") + "\n\n")

	// Enhanced help text using renderKeyHelp
	helpKeys := map[string]string{
		"↑/↓":   "Navigate",
		"Enter": "Select",
		"q":     "Quit",
	}
	helpText := renderKeyHelp(helpKeys)
	content.WriteString(helpText)

	return "\n" + borderStyleThick.Render(content.String())
}

// SSH Management Types
type sshItem struct {
	config SSHConfig
}

func (i sshItem) Title() string { return i.config.Name }
func (i sshItem) Description() string {
	return fmt.Sprintf("%s@%s:%d", i.config.User, i.config.Host, i.config.Port)
}
func (i sshItem) FilterValue() string { return i.config.Name }

// SSHListModel represents the SSH configuration list
type SSHListModel struct {
	list     list.Model
	config   *Config
	quitting bool
}

// NewSSHListModel creates a new SSH list model
func NewSSHListModel(config *Config) SSHListModel {
	items := make([]list.Item, len(config.SSHConfigs))
	for i, cfg := range config.SSHConfigs {
		items[i] = sshItem{config: cfg}
	}

	l := list.New(items, getCustomDelegate(), 0, 0)
	l.Title = "SSH Configurations"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SSHListModel{
		list:   l,
		config: config,
	}
}

func (m SSHListModel) Init() tea.Cmd {
	return nil
}

func (m SSHListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := titleStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Add new SSH config
			return NewSSHFormModel(m.config, -1), nil
		case "e":
			// Edit selected SSH config
			if i, ok := m.list.SelectedItem().(sshItem); ok {
				for idx, cfg := range m.config.SSHConfigs {
					if cfg.Name == i.config.Name {
						return NewSSHFormModel(m.config, idx), nil
					}
				}
			}
		case "d":
			// Delete selected SSH config
			if i, ok := m.list.SelectedItem().(sshItem); ok {
				for idx, cfg := range m.config.SSHConfigs {
					if cfg.Name == i.config.Name {
						m.config.SSHConfigs = append(m.config.SSHConfigs[:idx], m.config.SSHConfigs[idx+1:]...)
						SaveConfig("config.json", m.config)
						// Refresh list
						items := make([]list.Item, len(m.config.SSHConfigs))
						for i, cfg := range m.config.SSHConfigs {
							items[i] = sshItem{config: cfg}
						}
						m.list.SetItems(items)
						break
					}
				}
			}
		case "t":
			// Test connection
			if i, ok := m.list.SelectedItem().(sshItem); ok {
				return NewSSHTestModel(i.config, m.config), nil
			}
		case "backspace", "esc":
			return NewMainMenu(m.config), nil
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SSHListModel) View() string {
	if m.quitting {
		return successStyle.Render("Goodbye!\n")
	}

	// Check if list is empty
	if len(m.list.Items()) == 0 {
		emptyMessage := subtitleStyle.Render("No SSH Configurations Found") + "\n\n" +
			bodyStyle.Render("You haven't added any SSH server configurations yet.") + "\n" +
			bodyStyle.Render("Press ") + renderIcon("add", lipgloss.NewStyle().Foreground(accentColor)) +
			bodyStyle.Render(" 'a' to add your first SSH configuration.")

		helpKeys := map[string]string{
			"a":   "Add",
			"esc": "Back",
			"q":   "Quit",
		}
		helpText := renderKeyHelp(helpKeys)

		return "\n" + borderStyle.Render(
			titleStyle.Render("SSH Configurations")+"\n\n"+
				emptyMessage+"\n\n"+
				renderDivider(60, dividerStyle)+"\n\n"+
				helpText,
		)
	}

	// Enhanced list view with connection status indicators
	var listContent strings.Builder

	// Add title
	listContent.WriteString(titleStyle.Render("SSH Configurations") + "\n\n")

	// Render each SSH config as an enhanced card
	for i, item := range m.list.Items() {
		if sshItem, ok := item.(sshItem); ok {
			// Connection status indicator (colored dot)
			statusIndicator := renderIcon("circle", lipgloss.NewStyle().Foreground(mutedColor))

			// Create card content
			nameStyle := lipgloss.NewStyle().Foreground(textPrimary).Bold(true)
			detailStyle := lipgloss.NewStyle().Foreground(textSecondary)

			cardContent := statusIndicator + " " + nameStyle.Render(sshItem.config.Name) + "\n" +
				detailStyle.Render(fmt.Sprintf("  %s@%s:%d", sshItem.config.User, sshItem.config.Host, sshItem.config.Port)) + "\n" +
				detailStyle.Render(fmt.Sprintf("  Auth: %s", sshItem.config.AuthType))

			// Highlight selected item
			if i == m.list.Index() {
				selectedCardStyle := lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(accentColor).
					Background(backgroundLight).
					Padding(1, 2).
					MarginTop(0).
					MarginBottom(1)
				listContent.WriteString(selectedCardStyle.Render(cardContent) + "\n")
			} else {
				normalCardStyle := lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(mutedColor).
					Padding(1, 2).
					MarginTop(0).
					MarginBottom(1)
				listContent.WriteString(normalCardStyle.Render(cardContent) + "\n")
			}
		}
	}

	// Add divider before help text
	listContent.WriteString("\n" + renderDivider(60, dividerStyle) + "\n\n")

	// Enhanced help text using renderKeyHelp
	helpKeys := map[string]string{
		"a":   "Add",
		"e":   "Edit",
		"d":   "Delete",
		"t":   "Test Connection",
		"esc": "Back",
		"q":   "Quit",
	}
	helpText := renderKeyHelp(helpKeys)
	listContent.WriteString(helpText)

	return "\n" + borderStyle.Render(listContent.String())
}

// SSHFormModel represents the SSH configuration form
type SSHFormModel struct {
	config    *Config
	editIndex int
	form      []formField
	cursor    int
}

type formField struct {
	label     string
	value     string
	fieldType string // "text", "password", "key", "port", "auth"
	isValid   bool
	errorMsg  string
	icon      string
	multiline bool
}

// validateFormField validates a form field based on its type
func validateFormField(field *formField) {
	field.isValid = true
	field.errorMsg = ""

	switch field.fieldType {
	case "text":
		if strings.TrimSpace(field.value) == "" {
			field.isValid = false
			field.errorMsg = "This field is required"
		}
	case "port":
		port := parseInt(field.value)
		if port <= 0 || port > 65535 {
			field.isValid = false
			field.errorMsg = "Port must be between 1 and 65535"
		}
	case "auth":
		validTypes := []string{"password", "key", "agent"}
		valid := false
		for _, t := range validTypes {
			if field.value == t {
				valid = true
				break
			}
		}
		if !valid {
			field.isValid = false
			field.errorMsg = "Must be 'password', 'key', or 'agent'"
		}
	}
}

func NewSSHFormModel(config *Config, editIndex int) SSHFormModel {
	var fields []formField

	if editIndex >= 0 && editIndex < len(config.SSHConfigs) {
		cfg := config.SSHConfigs[editIndex]
		fields = []formField{
			{label: "Name", value: cfg.Name, fieldType: "text", isValid: true, icon: "edit"},
			{label: "Host", value: cfg.Host, fieldType: "text", isValid: true, icon: "server"},
			{label: "Port", value: fmt.Sprintf("%d", cfg.Port), fieldType: "port", isValid: true, icon: "gear"},
			{label: "User", value: cfg.User, fieldType: "text", isValid: true, icon: "user"},
			{label: "Auth Type", value: cfg.AuthType, fieldType: "auth", isValid: true, icon: "lock"},
			{label: "Password", value: cfg.Password, fieldType: "password", isValid: true, icon: "key"},
			{label: "Key File", value: cfg.KeyFile, fieldType: "key", isValid: true, icon: "file"},
			{label: "Key Password", value: cfg.KeyPass, fieldType: "password", isValid: true, icon: "lock"},
		}
	} else {
		fields = []formField{
			{label: "Name", value: "", fieldType: "text", isValid: true, icon: "edit"},
			{label: "Host", value: "", fieldType: "text", isValid: true, icon: "server"},
			{label: "Port", value: "22", fieldType: "port", isValid: true, icon: "gear"},
			{label: "User", value: "", fieldType: "text", isValid: true, icon: "user"},
			{label: "Auth Type", value: "password", fieldType: "auth", isValid: true, icon: "lock"},
			{label: "Password", value: "", fieldType: "password", isValid: true, icon: "key"},
			{label: "Key File", value: "", fieldType: "key", isValid: true, icon: "file"},
			{label: "Key Password", value: "", fieldType: "password", isValid: true, icon: "lock"},
		}
	}

	return SSHFormModel{
		config:    config,
		editIndex: editIndex,
		form:      fields,
	}
}

func (m SSHFormModel) Init() tea.Cmd {
	return nil
}

func (m SSHFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.form)-1 {
				m.cursor++
			}
		case "enter":
			// Validate current field before moving forward
			validateFormField(&m.form[m.cursor])

			if m.cursor == len(m.form)-1 {
				// Validate all fields before saving
				allValid := true
				for i := range m.form {
					validateFormField(&m.form[i])
					if !m.form[i].isValid {
						allValid = false
					}
				}

				if !allValid {
					// Don't save if validation fails, stay on form
					return m, nil
				}

				// Save the configuration
				cfg := SSHConfig{
					Name:     m.form[0].value,
					Host:     m.form[1].value,
					Port:     parseInt(m.form[2].value),
					User:     m.form[3].value,
					AuthType: m.form[4].value,
					Password: m.form[5].value,
					KeyFile:  m.form[6].value,
					KeyPass:  m.form[7].value,
				}

				if m.editIndex >= 0 {
					m.config.SSHConfigs[m.editIndex] = cfg
				} else {
					m.config.SSHConfigs = append(m.config.SSHConfigs, cfg)
				}

				SaveConfig("config.json", m.config)
				return NewSSHListModel(m.config), nil
			} else {
				// Move to next field if current field is valid
				if m.form[m.cursor].isValid {
					m.cursor++
				}
			}
		case "backspace":
			if len(m.form[m.cursor].value) > 0 {
				m.form[m.cursor].value = m.form[m.cursor].value[:len(m.form[m.cursor].value)-1]
				// Revalidate on change
				validateFormField(&m.form[m.cursor])
			}
		case "esc":
			return NewSSHListModel(m.config), nil
		default:
			if len(msg.String()) == 1 {
				m.form[m.cursor].value += msg.String()
				// Revalidate on change
				validateFormField(&m.form[m.cursor])
			}
		}
	}

	return m, nil
}

func (m SSHFormModel) View() string {
	var formContent strings.Builder

	// Determine form title
	formTitle := "Add SSH Configuration"
	if m.editIndex >= 0 {
		formTitle = "Edit SSH Configuration"
	}

	// Add title with decorative divider
	formContent.WriteString(titleStyle.Render(formTitle) + "\n\n")
	formContent.WriteString(renderDividerDecorative(70, dividerStyleAccent, "dots") + "\n\n")

	// Define styles for form rendering
	fieldLabelStyle := lipgloss.NewStyle().
		Foreground(textSecondary).
		Bold(true).
		Width(20).
		Align(lipgloss.Left)

	fieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		Padding(0, 1).
		Width(40)

	activeFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(0, 1).
		Width(40)

	errorFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(0, 1).
		Width(40)

	cursorIndicatorStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)

	// Render each form field
	for i, field := range m.form {
		// Field type indicator icon
		iconStyle := lipgloss.NewStyle().Foreground(primaryColor)
		fieldIcon := renderIcon(field.icon, iconStyle)

		// Cursor indicator for active field
		cursorIndicator := "  "
		if m.cursor == i {
			cursorIndicator = cursorIndicatorStyle.Render("▶ ")
		}

		// Label with icon
		labelText := fieldLabelStyle.Render(field.label)

		// Value display (mask password fields)
		displayValue := field.value
		if field.fieldType == "password" && field.value != "" {
			displayValue = strings.Repeat("•", len(field.value))
		}

		// Add cursor position indicator for active field
		if m.cursor == i {
			displayValue = displayValue + cursorIndicatorStyle.Render("│")
		}

		// Choose appropriate style based on field state
		var valueRendered string
		if m.cursor == i {
			valueRendered = activeFieldValueStyle.Render(displayValue)
		} else if !field.isValid {
			valueRendered = errorFieldValueStyle.Render(displayValue)
		} else {
			valueRendered = fieldValueStyle.Render(displayValue)
		}

		// Compose the field line
		fieldLine := cursorIndicator + fieldIcon + " " + labelText + " " + valueRendered

		formContent.WriteString(fieldLine + "\n")

		// Show error message if field is invalid
		if !field.isValid && field.errorMsg != "" {
			errorMsgStyle := lipgloss.NewStyle().
				Foreground(errorColor).
				Italic(true).
				MarginLeft(25)
			formContent.WriteString(errorMsgStyle.Render("  ⚠ "+field.errorMsg) + "\n")
		}

		// Add spacing between fields
		formContent.WriteString("\n")
	}

	// Add help text with decorative divider
	formContent.WriteString(renderDividerDecorative(70, dividerStyle, "default") + "\n\n")

	// Enhanced help text using renderKeyHelp
	helpKeys := map[string]string{
		"↑/↓":   "Navigate",
		"Enter": "Next/Save",
		"Esc":   "Cancel",
	}
	helpText := renderKeyHelp(helpKeys)
	formContent.WriteString(helpText + "\n")

	// Additional context help
	contextHelp := helpTextStyle.Render("Fill in all fields and press Enter on the last field to save")
	formContent.WriteString("\n" + contextHelp)

	return "\n" + formStyle.Render(formContent.String())
}

// SSHTestModel represents the SSH connection test
type SSHTestModel struct {
	config       SSHConfig
	result       string
	done         bool
	quitting     bool
	appConfig    *Config
	spinnerFrame int
}

func NewSSHTestModel(sshConfig SSHConfig, appConfig *Config) SSHTestModel {
	return SSHTestModel{
		config:    sshConfig,
		appConfig: appConfig,
	}
}

func (m SSHTestModel) Init() tea.Cmd {
	return tea.Batch(m.testConnection(), m.tickCmd())
}

func (m SSHTestModel) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (m SSHTestModel) testConnection() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		auth, err := m.config.GetAuthMethod()
		if err != nil {
			return testResultMsg{result: fmt.Sprintf("Auth error: %v", err)}
		}

		client, err := ssh.New(m.config.User, m.config.Host, auth)
		if err != nil {
			return testResultMsg{result: fmt.Sprintf("Connection failed: %v", err)}
		}
		defer client.Close()

		output, err := client.Run("echo 'SSH connection successful'")
		if err != nil {
			return testResultMsg{result: fmt.Sprintf("Command failed: %v", err)}
		}

		return testResultMsg{result: fmt.Sprintf("Success: %s", string(output))}
	})
}

type testResultMsg struct {
	result string
}

func (m SSHTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case testResultMsg:
		m.result = msg.result
		m.done = true
		return m, nil

	case tickMsg:
		if !m.done {
			m.spinnerFrame++
			return m, m.tickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		if m.done && msg.String() == "enter" {
			return NewSSHListModel(m.appConfig), nil
		}
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m SSHTestModel) View() string {
	var content strings.Builder

	// Title with decorative divider
	content.WriteString(titleStyle.Render(fmt.Sprintf("Testing Connection to %s", m.config.Name)) + "\n\n")
	content.WriteString(renderDividerDecorative(70, dividerStyleAccent, "dots") + "\n\n")

	if !m.done {
		// Show animated spinner during connection test
		spinner := renderSpinner(m.spinnerFrame)
		loadingIcon := renderIcon("loading", statusStyle)
		content.WriteString(spinner + " " + statusStyle.Render("Testing connection...") + " " + loadingIcon + "\n\n")

		// Show connection details in styled info box using renderCard
		connectionDetails := fmt.Sprintf("Host: %s\nPort: %d\nUser: %s\nAuth: %s",
			m.config.Host,
			m.config.Port,
			m.config.User,
			m.config.AuthType)
		content.WriteString(renderCard(connectionDetails, "Connection Details") + "\n\n")

		// Add divider before help text
		content.WriteString(renderDivider(70, dividerStyle) + "\n\n")

		// Enhanced help text using renderKeyHelp
		helpKeys := map[string]string{
			"q": "Quit",
		}
		helpText := renderKeyHelp(helpKeys)
		content.WriteString(helpText)
	} else {
		// Enhance result display with icons and color coding
		if strings.Contains(m.result, "Success") {
			// Success case with icon
			successIcon := renderIcon("success", successStyle)
			content.WriteString(successIcon + " " + successStyle.Render(m.result) + "\n\n")

			// Show connection details in success card
			connectionDetails := fmt.Sprintf("Host: %s\nPort: %d\nUser: %s\nAuth: %s\n\nStatus: Connected",
				m.config.Host,
				m.config.Port,
				m.config.User,
				m.config.AuthType)
			content.WriteString(renderCard(connectionDetails, "Connection Details") + "\n\n")
		} else {
			// Error case with icon and improved formatting
			errorIcon := renderIcon("error", errorStyle)
			content.WriteString(errorIcon + " " + errorStyle.Render("Connection Failed") + "\n\n")

			// Show error details in error-styled card
			errorDetails := fmt.Sprintf("Host: %s\nPort: %d\nUser: %s\nAuth: %s\n\nError:\n%s",
				m.config.Host,
				m.config.Port,
				m.config.User,
				m.config.AuthType,
				m.result)

			// Create error card style
			errorCardStyle := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(errorColor).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(1)

			errorCardContent := subtitleStyle.Render("Connection Details") + "\n\n" + errorDetails
			content.WriteString(errorCardStyle.Render(errorCardContent) + "\n\n")
		}

		// Add divider before help text
		content.WriteString(renderDivider(70, dividerStyle) + "\n\n")

		// Enhanced help text using renderKeyHelp
		helpKeys := map[string]string{
			"Enter": "Return to SSH List",
			"q":     "Quit",
		}
		helpText := renderKeyHelp(helpKeys)
		content.WriteString(helpText)
	}

	return "\n" + borderStyle.Render(content.String())
}

// Project Management Types
type projectItem struct {
	project Project
}

func (i projectItem) Title() string       { return i.project.Name }
func (i projectItem) Description() string { return fmt.Sprintf("Servers: %v", i.project.DeployServers) }
func (i projectItem) FilterValue() string { return i.project.Name }

// ProjectListModel represents the project list
type ProjectListModel struct {
	list     list.Model
	config   *Config
	quitting bool
}

// NewProjectListModel creates a new project list model
func NewProjectListModel(config *Config) ProjectListModel {
	items := make([]list.Item, len(config.Projects))
	for i, proj := range config.Projects {
		items[i] = projectItem{project: proj}
	}

	l := list.New(items, getCustomDelegate(), 0, 0)
	l.Title = "Projects"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return ProjectListModel{
		list:   l,
		config: config,
	}
}

func (m ProjectListModel) Init() tea.Cmd {
	return nil
}

func (m ProjectListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := titleStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Add new project
			return NewProjectFormModel(m.config, -1), nil
		case "e":
			// Edit selected project
			if i, ok := m.list.SelectedItem().(projectItem); ok {
				for idx, proj := range m.config.Projects {
					if proj.Name == i.project.Name {
						return NewProjectFormModel(m.config, idx), nil
					}
				}
			}
		case "d":
			// Delete selected project
			if i, ok := m.list.SelectedItem().(projectItem); ok {
				for idx, proj := range m.config.Projects {
					if proj.Name == i.project.Name {
						m.config.Projects = append(m.config.Projects[:idx], m.config.Projects[idx+1:]...)
						SaveConfig("config.json", m.config)
						// Refresh list
						items := make([]list.Item, len(m.config.Projects))
						for i, proj := range m.config.Projects {
							items[i] = projectItem{project: proj}
						}
						m.list.SetItems(items)
						break
					}
				}
			}
		case "enter":
			// Deploy selected project
			if i, ok := m.list.SelectedItem().(projectItem); ok {
				return NewDeployModel(i.project, m.config), nil
			}
		case "backspace", "esc":
			return NewMainMenu(m.config), nil
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ProjectListModel) View() string {
	if m.quitting {
		return successStyle.Render("Goodbye!\n")
	}

	// Check if list is empty
	if len(m.list.Items()) == 0 {
		emptyMessage := subtitleStyle.Render("No Projects Found") + "\n\n" +
			bodyStyle.Render("You haven't added any projects yet.") + "\n" +
			bodyStyle.Render("Press ") + renderIcon("add", lipgloss.NewStyle().Foreground(accentColor)) +
			bodyStyle.Render(" 'a' to add your first project.")

		helpKeys := map[string]string{
			"a":   "Add",
			"esc": "Back",
			"q":   "Quit",
		}
		helpText := renderKeyHelp(helpKeys)

		return "\n" + borderStyle.Render(
			titleStyle.Render("Projects")+"\n\n"+
				emptyMessage+"\n\n"+
				renderDivider(60, dividerStyle)+"\n\n"+
				helpText,
		)
	}

	// Enhanced list view with project status badges
	var listContent strings.Builder

	// Add title
	listContent.WriteString(titleStyle.Render("Projects") + "\n\n")

	// Render each project as an enhanced card
	for i, item := range m.list.Items() {
		if projItem, ok := item.(projectItem); ok {
			// Add project status badge using renderBadge
			statusBadge := renderBadge("Active", "success")

			// Create card content with improved visual hierarchy
			nameStyle := lipgloss.NewStyle().Foreground(textPrimary).Bold(true)
			detailStyle := lipgloss.NewStyle().Foreground(textSecondary)

			// Project name with status badge
			projectHeader := nameStyle.Render(projItem.project.Name) + " " + statusBadge

			// Enhance server list display with chips/tags
			var serverChips []string
			for _, server := range projItem.project.DeployServers {
				serverChip := renderBadge(server, "muted")
				serverChips = append(serverChips, serverChip)
			}
			serverDisplay := strings.Join(serverChips, " ")

			// Build card content with improved visual hierarchy
			var cardContent strings.Builder
			cardContent.WriteString(projectHeader + "\n\n")

			// Server list with icon
			serverIcon := renderIcon("server", lipgloss.NewStyle().Foreground(primaryColor))
			if len(projItem.project.DeployServers) > 0 {
				cardContent.WriteString(detailStyle.Render("  "+serverIcon+" Servers: ") + serverDisplay + "\n")
			} else {
				cardContent.WriteString(detailStyle.Render("  "+serverIcon+" Servers: None configured") + "\n")
			}

			// Build instructions preview (if available)
			if projItem.project.BuildInstructions != "" {
				buildIcon := renderIcon("wrench", lipgloss.NewStyle().Foreground(infoColor))
				buildPreview := projItem.project.BuildInstructions
				if len(buildPreview) > 50 {
					buildPreview = buildPreview[:47] + "..."
				}
				cardContent.WriteString(detailStyle.Render("  "+buildIcon+" Build: "+buildPreview) + "\n")
			}

			// Deploy script preview (if available)
			if projItem.project.DeployScript != "" {
				deployIcon := renderIcon("rocket", lipgloss.NewStyle().Foreground(secondaryColor))
				deployPreview := projItem.project.DeployScript
				if len(deployPreview) > 50 {
					deployPreview = deployPreview[:47] + "..."
				}
				cardContent.WriteString(detailStyle.Render("  "+deployIcon+" Deploy: "+deployPreview) + "\n")
			}

			// Highlight selected item
			if i == m.list.Index() {
				selectedCardStyle := lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(accentColor).
					Background(backgroundLight).
					Padding(1, 2).
					MarginTop(0).
					MarginBottom(1)
				listContent.WriteString(selectedCardStyle.Render(cardContent.String()) + "\n")
			} else {
				normalCardStyle := lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(mutedColor).
					Padding(1, 2).
					MarginTop(0).
					MarginBottom(1)
				listContent.WriteString(normalCardStyle.Render(cardContent.String()) + "\n")
			}
		}
	}

	// Add divider before help text
	listContent.WriteString("\n" + renderDivider(60, dividerStyle) + "\n\n")

	// Update help text formatting with renderKeyHelp
	helpKeys := map[string]string{
		"a":     "Add",
		"e":     "Edit",
		"d":     "Delete",
		"Enter": "Deploy",
		"esc":   "Back",
		"q":     "Quit",
	}
	helpText := renderKeyHelp(helpKeys)
	listContent.WriteString(helpText)

	return "\n" + borderStyle.Render(listContent.String())
}

// ProjectFormModel represents the project configuration form
type ProjectFormModel struct {
	config    *Config
	editIndex int
	form      []formField
	cursor    int
}

func NewProjectFormModel(config *Config, editIndex int) ProjectFormModel {
	var fields []formField

	if editIndex >= 0 && editIndex < len(config.Projects) {
		proj := config.Projects[editIndex]
		fields = []formField{
			{label: "Name", value: proj.Name, fieldType: "text", isValid: true, icon: "edit"},
			{label: "Build Instructions", value: proj.BuildInstructions, fieldType: "text", isValid: true, icon: "wrench", multiline: true},
			{label: "Deploy Script", value: proj.DeployScript, fieldType: "text", isValid: true, icon: "rocket", multiline: true},
			{label: "Deploy Servers", value: fmt.Sprintf("%v", proj.DeployServers), fieldType: "text", isValid: true, icon: "server"},
		}
	} else {
		fields = []formField{
			{label: "Name", value: "", fieldType: "text", isValid: true, icon: "edit"},
			{label: "Build Instructions", value: "", fieldType: "text", isValid: true, icon: "wrench", multiline: true},
			{label: "Deploy Script", value: "", fieldType: "text", isValid: true, icon: "rocket", multiline: true},
			{label: "Deploy Servers", value: "", fieldType: "text", isValid: true, icon: "server"},
		}
	}

	return ProjectFormModel{
		config:    config,
		editIndex: editIndex,
		form:      fields,
	}
}

func (m ProjectFormModel) Init() tea.Cmd {
	return nil
}

func (m ProjectFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.form)-1 {
				m.cursor++
			}
		case "enter":
			// Validate current field before moving forward
			validateFormField(&m.form[m.cursor])

			if m.cursor == len(m.form)-1 {
				// Validate all fields before saving
				allValid := true
				for i := range m.form {
					validateFormField(&m.form[i])
					if !m.form[i].isValid {
						allValid = false
					}
				}

				if !allValid {
					// Don't save if validation fails, stay on form
					return m, nil
				}

				// Save the project
				proj := Project{
					Name:              m.form[0].value,
					BuildInstructions: m.form[1].value,
					DeployScript:      m.form[2].value,
					DeployServers:     parseServers(m.form[3].value),
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}

				if m.editIndex >= 0 {
					proj.CreatedAt = m.config.Projects[m.editIndex].CreatedAt
					m.config.Projects[m.editIndex] = proj
				} else {
					m.config.Projects = append(m.config.Projects, proj)
				}

				SaveConfig("config.json", m.config)
				return NewProjectListModel(m.config), nil
			} else {
				// Move to next field if current field is valid
				if m.form[m.cursor].isValid {
					m.cursor++
				}
			}
		case "backspace":
			if len(m.form[m.cursor].value) > 0 {
				m.form[m.cursor].value = m.form[m.cursor].value[:len(m.form[m.cursor].value)-1]
				// Revalidate on change
				validateFormField(&m.form[m.cursor])
			}
		case "esc":
			return NewProjectListModel(m.config), nil
		default:
			if len(msg.String()) == 1 {
				m.form[m.cursor].value += msg.String()
				// Revalidate on change
				validateFormField(&m.form[m.cursor])
			}
		}
	}

	return m, nil
}

func (m ProjectFormModel) View() string {
	var formContent strings.Builder

	// Determine form title
	formTitle := "Add Project"
	if m.editIndex >= 0 {
		formTitle = "Edit Project"
	}

	// Add title with decorative divider
	formContent.WriteString(titleStyle.Render(formTitle) + "\n\n")
	formContent.WriteString(renderDividerDecorative(70, dividerStyleAccent, "dots") + "\n\n")

	// Define styles for form rendering
	fieldLabelStyle := lipgloss.NewStyle().
		Foreground(textSecondary).
		Bold(true).
		Width(20).
		Align(lipgloss.Left)

	fieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		Padding(0, 1).
		Width(40)

	activeFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(0, 1).
		Width(40)

	errorFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(0, 1).
		Width(40)

	cursorIndicatorStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)

	multilineFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		Padding(1, 1).
		Width(40).
		Height(3)

	activeMultilineFieldValueStyle := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(backgroundLight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 1).
		Width(40).
		Height(3)

	// Render each form field
	for i, field := range m.form {
		// Field type indicator icon
		iconStyle := lipgloss.NewStyle().Foreground(primaryColor)
		fieldIcon := renderIcon(field.icon, iconStyle)

		// Cursor indicator for active field
		cursorIndicator := "  "
		if m.cursor == i {
			cursorIndicator = cursorIndicatorStyle.Render("▶ ")
		}

		// Label with icon
		labelText := fieldLabelStyle.Render(field.label)

		// Value display
		displayValue := field.value

		// Add cursor position indicator for active field
		if m.cursor == i {
			displayValue = displayValue + cursorIndicatorStyle.Render("│")
		}

		// Choose appropriate style based on field state and type
		var valueRendered string
		if field.multiline {
			// Multi-line field support with better visualization
			multilineContent := displayValue
			if len(multilineContent) > 120 {
				// Truncate for display with ellipsis
				multilineContent = multilineContent[:117] + "..."
			}

			if m.cursor == i {
				valueRendered = activeMultilineFieldValueStyle.Render(multilineContent)
			} else if !field.isValid {
				errorMultilineStyle := lipgloss.NewStyle().
					Foreground(textPrimary).
					Background(backgroundLight).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(errorColor).
					Padding(1, 1).
					Width(40).
					Height(3)
				valueRendered = errorMultilineStyle.Render(multilineContent)
			} else {
				valueRendered = multilineFieldValueStyle.Render(multilineContent)
			}
		} else {
			// Single-line field
			if m.cursor == i {
				valueRendered = activeFieldValueStyle.Render(displayValue)
			} else if !field.isValid {
				valueRendered = errorFieldValueStyle.Render(displayValue)
			} else {
				valueRendered = fieldValueStyle.Render(displayValue)
			}
		}

		// Compose the field line
		fieldLine := cursorIndicator + fieldIcon + " " + labelText + " " + valueRendered

		formContent.WriteString(fieldLine + "\n")

		// Show error message if field is invalid
		if !field.isValid && field.errorMsg != "" {
			errorMsgStyle := lipgloss.NewStyle().
				Foreground(errorColor).
				Italic(true).
				MarginLeft(25)
			formContent.WriteString(errorMsgStyle.Render("  ⚠ "+field.errorMsg) + "\n")
		}

		// Add spacing between fields
		formContent.WriteString("\n")
	}

	// Add help text with decorative divider
	formContent.WriteString(renderDividerDecorative(70, dividerStyle, "default") + "\n\n")

	// Enhanced help text using renderKeyHelp
	helpKeys := map[string]string{
		"↑/↓":   "Navigate",
		"Enter": "Next/Save",
		"Esc":   "Cancel",
	}
	helpText := renderKeyHelp(helpKeys)
	formContent.WriteString(helpText + "\n")

	// Additional context help
	contextHelp := helpTextStyle.Render("Fill in all fields and press Enter on the last field to save")
	formContent.WriteString("\n" + contextHelp)

	return "\n" + formStyle.Render(formContent.String())
}

// DeployModel represents the deployment process
type DeployModel struct {
	project  Project
	config   *Config
	logs     []string
	done     bool
	quitting bool
}

func NewDeployModel(project Project, config *Config) DeployModel {
	return DeployModel{
		project: project,
		config:  config,
		logs:    []string{"Starting deployment..."},
	}
}

func (m DeployModel) Init() tea.Cmd {
	return m.startDeployment()
}

func (m DeployModel) startDeployment() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		logs := []string{"Starting deployment..."}

		// Find SSH configs for deploy servers
		for _, serverName := range m.project.DeployServers {
			var serverConfig *SSHConfig
			for _, cfg := range m.config.SSHConfigs {
				if cfg.Name == serverName {
					serverConfig = &cfg
					break
				}
			}

			if serverConfig == nil {
				logs = append(logs, fmt.Sprintf("Error: SSH config '%s' not found", serverName))
				continue
			}

			logs = append(logs, fmt.Sprintf("Connecting to %s...", serverName))

			auth, err := serverConfig.GetAuthMethod()
			if err != nil {
				logs = append(logs, fmt.Sprintf("Auth error for %s: %v", serverName, err))
				continue
			}

			client, err := ssh.New(serverConfig.User, serverConfig.Host, auth)
			if err != nil {
				logs = append(logs, fmt.Sprintf("Connection failed to %s: %v", serverName, err))
				continue
			}

			// Run build instructions if any
			if m.project.BuildInstructions != "" {
				logs = append(logs, fmt.Sprintf("Running build on %s...", serverName))
				output, err := client.Run(m.project.BuildInstructions)
				if err != nil {
					logs = append(logs, fmt.Sprintf("Build failed on %s: %v", serverName, err))
					client.Close()
					continue
				}
				logs = append(logs, fmt.Sprintf("Build output: %s", string(output)))
			}

			// Run deploy script
			logs = append(logs, fmt.Sprintf("Running deploy script on %s...", serverName))
			output, err := client.Run(m.project.DeployScript)
			if err != nil {
				logs = append(logs, fmt.Sprintf("Deploy failed on %s: %v", serverName, err))
			} else {
				logs = append(logs, fmt.Sprintf("Deploy successful on %s: %s", serverName, string(output)))
			}

			client.Close()
		}

		logs = append(logs, "Deployment finished.")
		return deployResultMsg{logs: logs}
	})
}

type deployResultMsg struct {
	logs []string
}

func (m DeployModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case deployResultMsg:
		m.logs = msg.logs
		m.done = true
		return m, nil

	case tea.KeyMsg:
		if m.done && msg.String() == "enter" {
			return NewProjectListModel(m.config), nil
		}
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m DeployModel) View() string {
	var content strings.Builder

	// Title with decorative divider
	content.WriteString(titleStyle.Render(fmt.Sprintf("Deploying %s", m.project.Name)) + "\n\n")
	content.WriteString(renderDividerDecorative(70, dividerStyleAccent, "arrows") + "\n\n")

	// Deployment progress visualization
	if !m.done {
		// Show progress bar during deployment
		totalSteps := len(m.project.DeployServers) * 2 // Connect + Deploy for each server
		if m.project.BuildInstructions != "" {
			totalSteps += len(m.project.DeployServers) // Add build step for each server
		}
		currentStep := len(m.logs) - 1 // Subtract 1 for "Starting deployment..." message
		if currentStep < 0 {
			currentStep = 0
		}
		if currentStep > totalSteps {
			currentStep = totalSteps
		}

		progressLabel := subtitleStyle.Render("Progress")
		content.WriteString(progressLabel + "\n")
		content.WriteString(renderProgressBar(currentStep, totalSteps, 60) + "\n\n")
	}

	// Enhanced log entries with timestamps, icons, and color-coded log levels
	logSectionTitle := subtitleStyle.Render("Deployment Log")
	content.WriteString(logSectionTitle + "\n")
	content.WriteString(renderDivider(70, dividerStyle) + "\n\n")

	// Determine log entry type and render with appropriate styling
	for i, log := range m.logs {
		var logIcon string
		var logStyle lipgloss.Style
		var logLevel string

		// Determine log level and styling based on content
		if strings.HasPrefix(log, "Error:") || strings.Contains(log, "failed") || strings.Contains(log, "Failed") {
			logIcon = renderIcon("error", errorStyle)
			logStyle = errorStyle
			logLevel = "ERROR"
		} else if strings.Contains(log, "successful") || strings.Contains(log, "Success") || strings.Contains(log, "finished") {
			logIcon = renderIcon("success", successStyle)
			logStyle = successStyle
			logLevel = "SUCCESS"
		} else if strings.Contains(log, "Starting") || strings.Contains(log, "Connecting") || strings.Contains(log, "Running") {
			logIcon = renderIcon("info", statusStyle)
			logStyle = statusStyle
			logLevel = "INFO"
		} else {
			logIcon = renderIcon("bullet", bodyStyle)
			logStyle = bodyStyle
			logLevel = "INFO"
		}

		// Add timestamp (simulated - using log index as time indicator)
		timestamp := fmt.Sprintf("[%02d:%02d]", i/60, i%60)
		timestampStyle := lipgloss.NewStyle().Foreground(textMuted)

		// Alternating log entry styles for readability
		var entryStyle lipgloss.Style
		if i%2 == 0 {
			entryStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginBottom(0)
		} else {
			entryStyle = lipgloss.NewStyle().
				Background(backgroundDark).
				Padding(0, 1).
				MarginBottom(0)
		}

		// Compose log entry with icon, timestamp, level badge, and message
		levelBadge := renderBadge(logLevel, strings.ToLower(logLevel))
		logEntry := logIcon + " " + timestampStyle.Render(timestamp) + " " + levelBadge + " " + logStyle.Render(log)

		content.WriteString(entryStyle.Render(logEntry) + "\n")
	}

	// Deployment summary section at the end (if done)
	if m.done {
		content.WriteString("\n" + renderDividerDecorative(70, dividerStyleAccent, "stars") + "\n\n")

		// Count successes and errors
		successCount := 0
		errorCount := 0
		for _, log := range m.logs {
			if strings.Contains(log, "successful") || strings.Contains(log, "Success") {
				successCount++
			} else if strings.HasPrefix(log, "Error:") || strings.Contains(log, "failed") || strings.Contains(log, "Failed") {
				errorCount++
			}
		}

		// Create summary card
		summaryTitle := subtitleStyle.Render("Deployment Summary")
		var summaryContent strings.Builder

		summaryContent.WriteString(summaryTitle + "\n\n")

		// Summary statistics
		totalServers := len(m.project.DeployServers)
		summaryContent.WriteString(bodyStyle.Render(fmt.Sprintf("Total Servers: %d", totalServers)) + "\n")
		summaryContent.WriteString(successStyle.Render(fmt.Sprintf("✓ Successful: %d", successCount)) + "\n")
		summaryContent.WriteString(errorStyle.Render(fmt.Sprintf("✗ Failed: %d", errorCount)) + "\n\n")

		// Overall status
		var overallStatus string
		var overallBadge string
		if errorCount == 0 {
			overallStatus = "Deployment completed successfully!"
			overallBadge = renderBadge("SUCCESS", "success")
		} else if successCount > 0 {
			overallStatus = "Deployment completed with some errors."
			overallBadge = renderBadge("PARTIAL", "warning")
		} else {
			overallStatus = "Deployment failed."
			overallBadge = renderBadge("FAILED", "error")
		}

		summaryContent.WriteString(bodyStyle.Render("Status: ") + overallBadge + " " + bodyStyle.Render(overallStatus) + "\n")

		content.WriteString(renderCard(summaryContent.String(), "") + "\n\n")

		// Help text
		helpKeys := map[string]string{
			"Enter": "Return to Projects",
			"q":     "Quit",
		}
		helpText := renderKeyHelp(helpKeys)
		content.WriteString(helpText)
	} else {
		// Show loading indicator during deployment
		content.WriteString("\n" + renderSpinner(len(m.logs)%10) + " " + statusStyle.Render("Deployment in progress...") + "\n\n")

		// Add divider before help text
		content.WriteString(renderDivider(70, dividerStyle) + "\n\n")

		// Enhanced help text using renderKeyHelp
		helpKeys := map[string]string{
			"q": "Quit",
		}
		helpText := renderKeyHelp(helpKeys)
		content.WriteString(helpText)
	}

	return borderStyle.Render(content.String())
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseServers(s string) []string {
	if s == "" {
		return []string{}
	}
	// Simple comma-separated parsing, could be improved
	var result []string
	current := ""
	for _, r := range s {
		if r == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if r != ' ' {
			current += string(r)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
