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

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 2).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#25A065")).
			Align(lipgloss.Center)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(lipgloss.Color("#F1F1F1"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#FFD700")).
				Bold(true).
				Background(lipgloss.Color("#333333"))

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4).
			Foreground(lipgloss.Color("#888888"))

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(lipgloss.Color("#AAAAAA")).
			Italic(true)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A5A5A")).
			Padding(1, 2)

	formStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#25A065")).
			Padding(1, 2).
			Width(80)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)
)

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
	return "\n" + borderStyle.Render(m.list.View())
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
	help := "• a: Add  • e: Edit  • d: Delete  • t: Test Connection  • esc: Back  • q: Quit"
	return "\n" + borderStyle.Render(m.list.View()+"\n\n"+help)
}

// SSHFormModel represents the SSH configuration form
type SSHFormModel struct {
	config    *Config
	editIndex int
	form      []field
	cursor    int
}

type field struct {
	label string
	value string
}

func NewSSHFormModel(config *Config, editIndex int) SSHFormModel {
	var fields []field
	var title string

	if editIndex >= 0 && editIndex < len(config.SSHConfigs) {
		cfg := config.SSHConfigs[editIndex]
		fields = []field{
			{"Name", cfg.Name},
			{"Host", cfg.Host},
			{"Port", fmt.Sprintf("%d", cfg.Port)},
			{"User", cfg.User},
			{"Auth Type (password/key/agent)", cfg.AuthType},
			{"Password", cfg.Password},
			{"Key File", cfg.KeyFile},
			{"Key Password", cfg.KeyPass},
		}
		_ = title // TODO: use title in view
	} else {
		fields = []field{
			{"Name", ""},
			{"Host", ""},
			{"Port", "22"},
			{"User", ""},
			{"Auth Type (password/key/agent)", "password"},
			{"Password", ""},
			{"Key File", ""},
			{"Key Password", ""},
		}
		title = "Add SSH Configuration"
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
			if m.cursor == len(m.form)-1 {
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
				m.cursor++
			}
		case "backspace":
			if len(m.form[m.cursor].value) > 0 {
				m.form[m.cursor].value = m.form[m.cursor].value[:len(m.form[m.cursor].value)-1]
			}
		case "esc":
			return NewSSHListModel(m.config), nil
		default:
			if len(msg.String()) == 1 {
				m.form[m.cursor].value += msg.String()
			}
		}
	}

	return m, nil
}

func (m SSHFormModel) View() string {
	var s string

	for i, field := range m.form {
		cursor := " "
		if m.cursor == i {
			cursor = statusStyle.Render(">")
		}
		s += fmt.Sprintf("%s %s: %s\n", cursor, field.label, field.value)
	}

	s += "\n" + helpStyle.Render("Press Enter to save, Esc to cancel")
	return formStyle.Render(titleStyle.Render("SSH Configuration Form") + "\n\n" + s)
}

// SSHTestModel represents the SSH connection test
type SSHTestModel struct {
	config    SSHConfig
	result    string
	done      bool
	quitting  bool
	appConfig *Config
}

func NewSSHTestModel(sshConfig SSHConfig, appConfig *Config) SSHTestModel {
	return SSHTestModel{
		config:    sshConfig,
		appConfig: appConfig,
	}
}

func (m SSHTestModel) Init() tea.Cmd {
	return m.testConnection()
}

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
	s := titleStyle.Render(fmt.Sprintf("Testing connection to %s", m.config.Name)) + "\n\n"

	if !m.done {
		s += statusStyle.Render("Testing connection...") + "\n"
	} else {
		if strings.Contains(m.result, "Success") {
			s += successStyle.Render(m.result) + "\n\n"
		} else {
			s += errorStyle.Render(m.result) + "\n\n"
		}
		s += helpStyle.Render("Press Enter to return")
	}

	return borderStyle.Render(s)
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
	help := "• a: Add  • e: Edit  • d: Delete  • enter: Deploy  • esc: Back  • q: Quit"
	return "\n" + borderStyle.Render(m.list.View()+"\n\n"+help)
}

// ProjectFormModel represents the project configuration form
type ProjectFormModel struct {
	config    *Config
	editIndex int
	form      []field
	cursor    int
}

func NewProjectFormModel(config *Config, editIndex int) ProjectFormModel {
	var fields []field

	if editIndex >= 0 && editIndex < len(config.Projects) {
		proj := config.Projects[editIndex]
		fields = []field{
			{"Name", proj.Name},
			{"Build Instructions", proj.BuildInstructions},
			{"Deploy Script", proj.DeployScript},
			{"Deploy Servers (comma-separated)", fmt.Sprintf("%v", proj.DeployServers)},
		}
	} else {
		fields = []field{
			{"Name", ""},
			{"Build Instructions", ""},
			{"Deploy Script", ""},
			{"Deploy Servers (comma-separated)", ""},
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
			if m.cursor == len(m.form)-1 {
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
				m.cursor++
			}
		case "backspace":
			if len(m.form[m.cursor].value) > 0 {
				m.form[m.cursor].value = m.form[m.cursor].value[:len(m.form[m.cursor].value)-1]
			}
		case "esc":
			return NewProjectListModel(m.config), nil
		default:
			if len(msg.String()) == 1 {
				m.form[m.cursor].value += msg.String()
			}
		}
	}

	return m, nil
}

func (m ProjectFormModel) View() string {
	var s string

	for i, field := range m.form {
		cursor := " "
		if m.cursor == i {
			cursor = statusStyle.Render(">")
		}
		s += fmt.Sprintf("%s %s: %s\n", cursor, field.label, field.value)
	}

	s += "\n" + helpStyle.Render("Press Enter to save, Esc to cancel")
	return formStyle.Render(titleStyle.Render("Project Configuration Form") + "\n\n" + s)
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
	s := titleStyle.Render(fmt.Sprintf("Deploying %s", m.project.Name)) + "\n\n"

	for _, log := range m.logs {
		if strings.HasPrefix(log, "Error:") || strings.Contains(log, "failed") || strings.Contains(log, "Failed") {
			s += errorStyle.Render(log) + "\n"
		} else if strings.Contains(log, "successful") || strings.Contains(log, "Success") || strings.Contains(log, "finished") {
			s += successStyle.Render(log) + "\n"
		} else {
			s += statusStyle.Render(log) + "\n"
		}
	}

	if m.done {
		s += "\n" + helpStyle.Render("Press Enter to return")
	}

	return borderStyle.Render(s)
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
