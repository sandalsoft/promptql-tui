package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sandalsoft/promptql-tui/internal/config"
	"github.com/sandalsoft/promptql-tui/internal/sdk"
)

type view int

const (
	viewSetup view = iota
	viewProjects
	viewThreads
	viewChat
)

type setupField int

const (
	fieldPAT setupField = iota
	fieldAPIKey
	fieldDDNURL
	fieldTimezone
	fieldCount
)

// ChatMessage represents a message displayed in the chat view.
type ChatMessage struct {
	Role    string // "user" or "assistant"
	Content string
}

// Model is the root Bubble Tea model for the TUI.
type Model struct {
	// State
	cfg     *config.Config
	client  *sdk.Client
	view    view
	width   int
	height  int
	err     error
	loading bool
	spinner spinner.Model

	// Setup view
	setupInputs []textinput.Model
	setupCursor int

	// Projects view
	projects       []sdk.UserProject
	projectCursor  int
	selectedProject *sdk.UserProject
	buildFQDN      string

	// Threads view
	threads      []sdk.Thread
	threadCursor int
	activeThread *sdk.Thread

	// Chat view
	chatInput textarea.Model
	messages  []ChatMessage
	threadID  string
}

// New creates a new TUI model.
func New(cfg *config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	// Setup inputs
	inputs := make([]textinput.Model, fieldCount)

	pat := textinput.New()
	pat.Placeholder = "your-personal-access-token"
	pat.CharLimit = 256
	pat.EchoMode = textinput.EchoPassword
	pat.EchoCharacter = '*'
	inputs[fieldPAT] = pat

	apiKey := textinput.New()
	apiKey.Placeholder = "your-api-key (optional, for direct query)"
	apiKey.CharLimit = 256
	inputs[fieldAPIKey] = apiKey

	ddnURL := textinput.New()
	ddnURL.Placeholder = "https://your-project.ddn.hasura.app/graphql"
	ddnURL.CharLimit = 512
	inputs[fieldDDNURL] = ddnURL

	tz := textinput.New()
	tz.Placeholder = "UTC"
	tz.CharLimit = 64
	inputs[fieldTimezone] = tz

	// Pre-fill from config
	if cfg.PAT != "" {
		inputs[fieldPAT].SetValue(cfg.PAT)
	}
	if cfg.APIKey != "" {
		inputs[fieldAPIKey].SetValue(cfg.APIKey)
	}
	if cfg.DDNURL != "" {
		inputs[fieldDDNURL].SetValue(cfg.DDNURL)
	}
	if cfg.Timezone != "" {
		inputs[fieldTimezone].SetValue(cfg.Timezone)
	}

	// Chat input
	ta := textarea.New()
	ta.Placeholder = "Ask PromptQL a question..."
	ta.CharLimit = 4096
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	m := Model{
		cfg:         cfg,
		spinner:     s,
		setupInputs: inputs,
		chatInput:   ta,
	}

	// Skip setup if already configured
	if cfg.HasCredentials() {
		m.view = viewProjects
		m.client = sdk.NewClient(sdk.ClientOptions{
			PAT:       cfg.PAT,
			APIKey:    cfg.APIKey,
			ProjectID: cfg.ProjectID,
		})
		m.loading = true
	} else {
		m.view = viewSetup
		m.setupInputs[0].Focus()
	}

	return m
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}
	if m.loading && m.view == viewProjects {
		cmds = append(cmds, m.loadProjects())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.handleEsc()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.chatInput.SetWidth(msg.Width - 4)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	switch m.view {
	case viewSetup:
		return m.updateSetup(msg)
	case viewProjects:
		return m.updateProjects(msg)
	case viewThreads:
		return m.updateThreads(msg)
	case viewChat:
		return m.updateChat(msg)
	}

	return m, nil
}

func (m Model) View() string {
	var content string
	switch m.view {
	case viewSetup:
		content = m.viewSetup()
	case viewProjects:
		content = m.viewProjects()
	case viewThreads:
		content = m.viewThreads()
	case viewChat:
		content = m.viewChat()
	}

	return content
}

// --- Setup View ---

func (m Model) viewSetup() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("PromptQL TUI Setup"))
	b.WriteString("\n\n")

	labels := []string{"PAT (Personal Access Token)", "API Key", "DDN URL", "Timezone"}
	for i, label := range labels {
		if i == m.setupCursor {
			b.WriteString(promptStyle.Render("> " + label))
		} else {
			b.WriteString(helpStyle.Render("  " + label))
		}
		b.WriteString("\n")
		b.WriteString("  " + m.setupInputs[i].View())
		b.WriteString("\n\n")
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("tab/shift+tab: navigate  |  enter: save & continue  |  ctrl+c: quit"))
	return b.String()
}

func (m Model) updateSetup(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.setupInputs[m.setupCursor].Blur()
			m.setupCursor = (m.setupCursor + 1) % int(fieldCount)
			m.setupInputs[m.setupCursor].Focus()
			return m, nil
		case "shift+tab", "up":
			m.setupInputs[m.setupCursor].Blur()
			m.setupCursor = (m.setupCursor - 1 + int(fieldCount)) % int(fieldCount)
			m.setupInputs[m.setupCursor].Focus()
			return m, nil
		case "enter":
			return m.saveSetup()
		}

	case configSavedMsg:
		m.view = viewProjects
		m.loading = true
		m.err = nil
		return m, tea.Batch(m.spinner.Tick, m.loadProjects())
	}

	var cmd tea.Cmd
	m.setupInputs[m.setupCursor], cmd = m.setupInputs[m.setupCursor].Update(msg)
	return m, cmd
}

func (m Model) saveSetup() (tea.Model, tea.Cmd) {
	pat := m.setupInputs[fieldPAT].Value()
	if pat == "" {
		m.err = fmt.Errorf("PAT is required")
		return m, nil
	}

	m.cfg.PAT = pat
	m.cfg.APIKey = m.setupInputs[fieldAPIKey].Value()
	m.cfg.DDNURL = m.setupInputs[fieldDDNURL].Value()
	m.cfg.Timezone = m.setupInputs[fieldTimezone].Value()
	if m.cfg.Timezone == "" {
		m.cfg.Timezone = "UTC"
	}

	m.client = sdk.NewClient(sdk.ClientOptions{
		PAT:    m.cfg.PAT,
		APIKey: m.cfg.APIKey,
	})

	return m, func() tea.Msg {
		if err := m.cfg.Save(); err != nil {
			return errMsg{err}
		}
		return configSavedMsg{}
	}
}

// --- Projects View ---

func (m Model) viewProjects() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("PromptQL Projects"))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(m.spinner.View() + " Loading projects...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: "+m.err.Error()) + "\n\n")
		b.WriteString(helpStyle.Render("r: retry  |  s: setup  |  ctrl+c: quit"))
		return b.String()
	}

	if len(m.projects) == 0 {
		b.WriteString(helpStyle.Render("No projects found.") + "\n\n")
		b.WriteString(helpStyle.Render("s: setup  |  ctrl+c: quit"))
		return b.String()
	}

	b.WriteString(subtitleStyle.Render(fmt.Sprintf("%d projects found", len(m.projects))))
	b.WriteString("\n\n")

	// Calculate visible window: header takes ~4 lines, footer takes ~2 lines
	maxVisible := m.height - 6
	if maxVisible < 3 {
		maxVisible = 3
	}
	if maxVisible > len(m.projects) {
		maxVisible = len(m.projects)
	}

	// Determine scroll offset to keep cursor visible
	scrollOffset := 0
	if m.projectCursor >= maxVisible {
		scrollOffset = m.projectCursor - maxVisible + 1
	}
	end := scrollOffset + maxVisible
	if end > len(m.projects) {
		end = len(m.projects)
		scrollOffset = end - maxVisible
		if scrollOffset < 0 {
			scrollOffset = 0
		}
	}

	for i := scrollOffset; i < end; i++ {
		p := m.projects[i]
		cursor := "  "
		style := normalItemStyle
		if i == m.projectCursor {
			cursor = "> "
			style = selectedItemStyle
		}
		line := fmt.Sprintf("%s%s", cursor, p.Name)
		if p.BuildFQDN != "" {
			line += helpStyle.Render(fmt.Sprintf("  (%s)", p.BuildFQDN))
		}
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	if end < len(m.projects) {
		b.WriteString(helpStyle.Render(fmt.Sprintf("  ... %d more below", len(m.projects)-end)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate  |  enter: select  |  s: setup  |  ctrl+c: quit"))
	return b.String()
}

func (m Model) updateProjects(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case projectsLoadedMsg:
		m.projects = msg.projects
		m.loading = false
		m.err = nil
		return m, nil

	case lookupResultMsg:
		m.buildFQDN = msg.result.BuildFQDN
		// Populate project ID from the lookup result
		if m.selectedProject != nil {
			m.selectedProject.ProjectID = msg.result.ProjectID
			m.selectedProject.BuildFQDN = msg.result.BuildFQDN
		}
		m.cfg.ProjectID = msg.result.ProjectID
		m.loading = false
		// Now load threads
		m.view = viewThreads
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.loadThreads())

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "j", "down":
			if m.projectCursor < len(m.projects)-1 {
				m.projectCursor++
			}
			return m, nil
		case "k", "up":
			if m.projectCursor > 0 {
				m.projectCursor--
			}
			return m, nil
		case "enter":
			if len(m.projects) > 0 {
				return m.selectProject()
			}
			return m, nil
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.loadProjects())
		case "s":
			m.view = viewSetup
			m.setupInputs[0].Focus()
			return m, nil
		}
	}

	return m, nil
}

func (m Model) selectProject() (tea.Model, tea.Cmd) {
	p := m.projects[m.projectCursor]
	m.selectedProject = &p
	m.loading = true
	m.err = nil

	return m, tea.Batch(m.spinner.Tick, m.lookupProjectByName(p.Name, p.BuildFQDN))
}

// --- Threads View ---

func (m Model) viewThreads() string {
	var b strings.Builder

	projectName := ""
	if m.selectedProject != nil {
		projectName = m.selectedProject.Name
	}
	b.WriteString(titleStyle.Render("Threads"))
	if projectName != "" {
		b.WriteString("  " + subtitleStyle.Render(projectName))
	}
	b.WriteString("\n")

	if m.loading {
		b.WriteString(m.spinner.View() + " Loading threads...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: "+m.err.Error()) + "\n\n")
		b.WriteString(helpStyle.Render("r: retry  |  esc: back  |  n: new thread  |  ctrl+c: quit"))
		return b.String()
	}

	b.WriteString(subtitleStyle.Render("Select a thread or start a new one"))
	b.WriteString("\n\n")

	// "New Thread" option always first
	if m.threadCursor == 0 {
		b.WriteString(selectedItemStyle.Render("> + New Thread"))
	} else {
		b.WriteString(normalItemStyle.Render("  + New Thread"))
	}
	b.WriteString("\n")

	for i, t := range m.threads {
		cursor := "  "
		style := normalItemStyle
		if i+1 == m.threadCursor {
			cursor = "> "
			style = selectedItemStyle
		}
		title := t.Title
		if title == "" {
			title = t.ThreadID
		}
		ts := ""
		if t.UpdatedAt != "" {
			ts = helpStyle.Render(fmt.Sprintf("  (%s)", t.UpdatedAt[:min(19, len(t.UpdatedAt))]))
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, title)) + ts)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate  |  enter: select  |  n: new thread  |  esc: back  |  ctrl+c: quit"))
	return b.String()
}

func (m Model) updateThreads(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case threadsLoadedMsg:
		m.threads = msg.threads
		m.loading = false
		m.err = nil
		return m, nil

	case threadStartedMsg:
		m.loading = false
		m.threadID = msg.result.ThreadID
		m.activeThread = &sdk.Thread{
			ThreadID: msg.result.ThreadID,
			Title:    msg.result.Title,
		}
		m.view = viewChat
		m.chatInput.Focus()
		m.messages = []ChatMessage{}

		// Add initial messages from thread events
		for _, evt := range msg.result.ThreadEvents {
			if content := extractEventContent(evt); content != "" {
				m.messages = append(m.messages, ChatMessage{
					Role:    "assistant",
					Content: content,
				})
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		maxIdx := len(m.threads) // 0 = new thread, 1..N = threads
		switch msg.String() {
		case "j", "down":
			if m.threadCursor < maxIdx {
				m.threadCursor++
			}
			return m, nil
		case "k", "up":
			if m.threadCursor > 0 {
				m.threadCursor--
			}
			return m, nil
		case "enter":
			if m.threadCursor == 0 {
				// New thread - go to chat
				m.view = viewChat
				m.threadID = ""
				m.activeThread = nil
				m.messages = []ChatMessage{}
				m.chatInput.Focus()
				return m, nil
			}
			// Resume existing thread
			return m.resumeThread(m.threads[m.threadCursor-1])
		case "n":
			m.view = viewChat
			m.threadID = ""
			m.activeThread = nil
			m.messages = []ChatMessage{}
			m.chatInput.Focus()
			return m, nil
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.loadThreads())
		}
	}

	return m, nil
}

func (m Model) resumeThread(t sdk.Thread) (tea.Model, tea.Cmd) {
	m.activeThread = &t
	m.threadID = t.ThreadID
	m.view = viewChat
	m.loading = true
	m.messages = []ChatMessage{}
	m.chatInput.Focus()

	return m, tea.Batch(m.spinner.Tick, m.loadEvents(t.ThreadID))
}

// --- Chat View ---

func (m Model) viewChat() string {
	var b strings.Builder

	threadTitle := "New Thread"
	if m.activeThread != nil && m.activeThread.Title != "" {
		threadTitle = m.activeThread.Title
	}
	b.WriteString(titleStyle.Render("Chat"))
	b.WriteString("  " + subtitleStyle.Render(threadTitle))
	b.WriteString("\n\n")

	if m.loading && len(m.messages) == 0 {
		b.WriteString(m.spinner.View() + " Loading...")
		return b.String()
	}

	// Messages
	maxMsgHeight := m.height - 12
	if maxMsgHeight < 5 {
		maxMsgHeight = 5
	}
	msgLines := []string{}
	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			msgLines = append(msgLines, userMsgStyle.Render("You: ")+msg.Content)
		case "assistant":
			msgLines = append(msgLines, assistantMsgStyle.Render("PromptQL: ")+msg.Content)
		}
		msgLines = append(msgLines, "")
	}

	// Scroll to show recent messages
	if len(msgLines) > maxMsgHeight {
		msgLines = msgLines[len(msgLines)-maxMsgHeight:]
	}

	b.WriteString(strings.Join(msgLines, "\n"))

	if m.loading {
		b.WriteString("\n" + m.spinner.View() + " Thinking...")
	}

	if m.err != nil {
		b.WriteString("\n" + errorStyle.Render("Error: "+m.err.Error()))
	}

	b.WriteString("\n\n")
	b.WriteString(promptStyle.Render("Message: "))
	b.WriteString("\n")
	b.WriteString(m.chatInput.View())
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("ctrl+s: send  |  esc: back to threads  |  ctrl+c: quit"))
	return b.String()
}

func (m Model) updateChat(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case eventsLoadedMsg:
		m.loading = false
		m.err = nil
		for _, evt := range msg.events {
			if content := extractEventContent(evt); content != "" {
				role := "assistant"
				if isUserEvent(evt) {
					role = "user"
				}
				m.messages = append(m.messages, ChatMessage{
					Role:    role,
					Content: content,
				})
			}
		}
		return m, nil

	case threadStartedMsg:
		m.loading = false
		m.err = nil
		m.threadID = msg.result.ThreadID
		m.activeThread = &sdk.Thread{
			ThreadID: msg.result.ThreadID,
			Title:    msg.result.Title,
		}
		for _, evt := range msg.result.ThreadEvents {
			if content := extractEventContent(evt); content != "" {
				m.messages = append(m.messages, ChatMessage{
					Role:    "assistant",
					Content: content,
				})
			}
		}
		return m, nil

	case messageSentMsg:
		m.loading = false
		m.err = nil
		if content := extractSendMessageContent(msg.result); content != "" {
			m.messages = append(m.messages, ChatMessage{
				Role:    "assistant",
				Content: content,
			})
		}
		return m, nil

	case queryResultMsg:
		m.loading = false
		m.err = nil
		if content := extractQueryResult(msg.result); content != "" {
			m.messages = append(m.messages, ChatMessage{
				Role:    "assistant",
				Content: content,
			})
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+s":
			return m.sendMessage()
		}
	}

	var cmd tea.Cmd
	m.chatInput, cmd = m.chatInput.Update(msg)
	return m, cmd
}

func (m Model) sendMessage() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(m.chatInput.Value())
	if text == "" {
		return m, nil
	}

	m.messages = append(m.messages, ChatMessage{
		Role:    "user",
		Content: text,
	})
	m.chatInput.Reset()
	m.loading = true
	m.err = nil

	// If we have a PAT and project selected, use threads API
	if m.client != nil && m.selectedProject != nil {
		if m.threadID == "" {
			// Start new thread
			return m, tea.Batch(m.spinner.Tick, m.startThread(text))
		}
		// Send to existing thread
		return m, tea.Batch(m.spinner.Tick, m.sendThreadMessage(text))
	}

	// If we only have an API key, use the query endpoint
	if m.cfg.APIKey != "" && m.cfg.DDNURL != "" {
		return m, tea.Batch(m.spinner.Tick, m.executeQuery(text))
	}

	m.err = fmt.Errorf("no project selected or API key + DDN URL configured")
	m.loading = false
	return m, nil
}

func (m Model) handleEsc() (tea.Model, tea.Cmd) {
	switch m.view {
	case viewChat:
		m.view = viewThreads
		m.chatInput.Blur()
		m.err = nil
		return m, nil
	case viewThreads:
		m.view = viewProjects
		m.err = nil
		return m, nil
	case viewProjects:
		m.view = viewSetup
		m.setupInputs[0].Focus()
		m.err = nil
		return m, nil
	}
	return m, nil
}

// --- Commands ---

func (m Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.client.Projects().ListUserProjects()
		if err != nil {
			return errMsg{err}
		}
		return projectsLoadedMsg{projects}
	}
}

func (m Model) lookupProject(projectID string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.Projects().Lookup(sdk.LookupOptions{
			ProjectID: projectID,
		})
		if err != nil {
			return errMsg{err}
		}
		return lookupResultMsg{result}
	}
}

func (m Model) lookupProjectByName(projectName string, fqdn string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.Projects().Lookup(sdk.LookupOptions{
			ProjectName: projectName,
			FQDN:        fqdn,
		})
		if err != nil {
			return errMsg{err}
		}
		return lookupResultMsg{result}
	}
}

func (m Model) loadThreads() tea.Cmd {
	return func() tea.Msg {
		if m.selectedProject == nil {
			return errMsg{fmt.Errorf("no project selected")}
		}
		threads, err := m.client.Threads().List(m.selectedProject.ProjectID, "")
		if err != nil {
			return errMsg{err}
		}
		return threadsLoadedMsg{threads}
	}
}

func (m Model) loadEvents(threadID string) tea.Cmd {
	return func() tea.Msg {
		events, err := m.client.Threads().GetEvents(threadID)
		if err != nil {
			return errMsg{err}
		}
		return eventsLoadedMsg{events}
	}
}

func (m Model) startThread(message string) tea.Cmd {
	return func() tea.Msg {
		if m.selectedProject == nil {
			return errMsg{fmt.Errorf("no project selected")}
		}
		result, err := m.client.Threads().Start(sdk.StartOptions{
			ProjectID: m.selectedProject.ProjectID,
			Message:   message,
			BuildFQDN: m.buildFQDN,
			Timezone:  m.cfg.Timezone,
		})
		if err != nil {
			return errMsg{err}
		}
		return threadStartedMsg{result}
	}
}

func (m Model) sendThreadMessage(message string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.Threads().SendMessage(sdk.SendMessageOptions{
			ThreadID:  m.threadID,
			Message:   message,
			BuildFQDN: m.buildFQDN,
			Timezone:  m.cfg.Timezone,
		})
		if err != nil {
			return errMsg{err}
		}
		return messageSentMsg{result}
	}
}

func (m Model) executeQuery(question string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.Query().Ask(
			question,
			m.cfg.DDNURL,
			nil,
			m.cfg.Timezone,
		)
		if err != nil {
			return errMsg{err}
		}
		return queryResultMsg{result}
	}
}

// --- Helpers ---

func extractEventContent(evt sdk.ThreadEvent) string {
	if evt.EventData == nil {
		return ""
	}
	// Try common event data shapes
	if text, ok := evt.EventData["assistant_message"].(map[string]interface{}); ok {
		if content, ok := text["text"].(string); ok {
			return content
		}
	}
	if text, ok := evt.EventData["user_message"].(map[string]interface{}); ok {
		if content, ok := text["text"].(string); ok {
			return content
		}
	}
	if text, ok := evt.EventData["message"].(string); ok {
		return text
	}
	if text, ok := evt.EventData["text"].(string); ok {
		return text
	}
	return ""
}

func isUserEvent(evt sdk.ThreadEvent) bool {
	if evt.EventData == nil {
		return false
	}
	if _, ok := evt.EventData["user_message"]; ok {
		return true
	}
	if role, ok := evt.EventData["role"].(string); ok {
		return role == "user"
	}
	return false
}

func extractSendMessageContent(result *sdk.SendMessageResult) string {
	if result == nil || result.EventData == nil {
		return ""
	}
	if text, ok := result.EventData["assistant_message"].(map[string]interface{}); ok {
		if content, ok := text["text"].(string); ok {
			return content
		}
	}
	if text, ok := result.EventData["message"].(string); ok {
		return text
	}
	if text, ok := result.EventData["text"].(string); ok {
		return text
	}
	return ""
}

func extractQueryResult(result map[string]interface{}) string {
	if result == nil {
		return ""
	}
	// Try to find the assistant response in common shapes
	if interactions, ok := result["interactions"].([]interface{}); ok {
		for i := len(interactions) - 1; i >= 0; i-- {
			if inter, ok := interactions[i].(map[string]interface{}); ok {
				if role, ok := inter["role"].(string); ok && role == "assistant" {
					if msg, ok := inter["assistant_message"].(map[string]interface{}); ok {
						if text, ok := msg["text"].(string); ok {
							return text
						}
					}
				}
			}
		}
	}
	if msg, ok := result["message"].(string); ok {
		return msg
	}
	if text, ok := result["text"].(string); ok {
		return text
	}
	// Last resort: stringify the whole response
	return fmt.Sprintf("%v", result)
}
