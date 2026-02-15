package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sandalsoft/promptql-tui/internal/config"
	"github.com/sandalsoft/promptql-tui/internal/sdk"
)

// ---------------------------------------------------------------------------
// View State Transitions
// ---------------------------------------------------------------------------

func TestNew_WithCredentials(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	if m.view != viewProjects {
		t.Errorf("expected viewProjects, got %d", m.view)
	}
	if !m.loading {
		t.Error("expected loading=true when credentials are present")
	}
	if m.client == nil {
		t.Error("expected SDK client to be initialized")
	}
}

func TestNew_WithoutCredentials(t *testing.T) {
	cfg := &config.Config{}
	m := New(cfg)
	if m.view != viewSetup {
		t.Errorf("expected viewSetup, got %d", m.view)
	}
	if m.loading {
		t.Error("expected loading=false when no credentials")
	}
}

func TestEsc_FromProjects(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.loading = false // simulate loading complete

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model := updated.(Model)
	if model.view != viewSetup {
		t.Errorf("expected viewSetup after esc from projects, got %d", model.view)
	}
}

func TestEsc_FromThreads(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewThreads
	m.loading = false

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model := updated.(Model)
	if model.view != viewProjects {
		t.Errorf("expected viewProjects after esc from threads, got %d", model.view)
	}
}

func TestEsc_FromChat(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewChat
	m.loading = false

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model := updated.(Model)
	if model.view != viewThreads {
		t.Errorf("expected viewThreads after esc from chat, got %d", model.view)
	}
}

// ---------------------------------------------------------------------------
// Project List Behavior
// ---------------------------------------------------------------------------

func TestProjectsLoaded(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.loading = true

	projects := []sdk.UserProject{
		{Name: "Alpha", DDNProjectID: "p-1", BuildFQDN: "alpha.app"},
		{Name: "Beta", DDNProjectID: "p-2"},
	}

	updated, _ := m.Update(projectsLoadedMsg{projects: projects})
	model := updated.(Model)

	if model.loading {
		t.Error("expected loading=false after projectsLoadedMsg")
	}
	if len(model.projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(model.projects))
	}
	if model.err != nil {
		t.Errorf("expected nil error, got %v", model.err)
	}
}

func TestProjectNavigation(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.loading = false
	m.projects = []sdk.UserProject{
		{Name: "A"}, {Name: "B"}, {Name: "C"},
	}
	m.projectCursor = 0

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := updated.(Model)
	if model.projectCursor != 1 {
		t.Errorf("after j: expected cursor=1, got %d", model.projectCursor)
	}

	// Move down again
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(Model)
	if model.projectCursor != 2 {
		t.Errorf("after j: expected cursor=2, got %d", model.projectCursor)
	}

	// Bounds check: can't go past last
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(Model)
	if model.projectCursor != 2 {
		t.Errorf("after j at end: expected cursor=2, got %d", model.projectCursor)
	}

	// Move up
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.projectCursor != 1 {
		t.Errorf("after k: expected cursor=1, got %d", model.projectCursor)
	}

	// Back to top
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.projectCursor != 0 {
		t.Errorf("after k: expected cursor=0, got %d", model.projectCursor)
	}

	// Bounds check: can't go above 0
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.projectCursor != 0 {
		t.Errorf("after k at top: expected cursor=0, got %d", model.projectCursor)
	}
}

func TestProjectsView_Error(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.loading = true

	updated, _ := m.Update(errMsg{err: errors.New("something went wrong")})
	model := updated.(Model)

	if model.loading {
		t.Error("expected loading=false after errMsg")
	}
	if model.err == nil {
		t.Fatal("expected error to be set")
	}

	output := model.View()
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("expected View output to contain error message, got: %s", output)
	}
}

func TestProjectsView_Empty(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.loading = false
	m.projects = []sdk.UserProject{}

	output := m.View()
	if !strings.Contains(output, "No projects found") {
		t.Errorf("expected View output to contain 'No projects found', got: %s", output)
	}
}

// ---------------------------------------------------------------------------
// Thread List Behavior
// ---------------------------------------------------------------------------

func TestThreadsLoaded(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewThreads
	m.loading = true

	threads := []sdk.Thread{
		{ThreadID: "t-1", Title: "Thread One"},
		{ThreadID: "t-2", Title: "Thread Two"},
	}

	updated, _ := m.Update(threadsLoadedMsg{threads: threads})
	model := updated.(Model)

	if model.loading {
		t.Error("expected loading=false after threadsLoadedMsg")
	}
	if len(model.threads) != 2 {
		t.Errorf("expected 2 threads, got %d", len(model.threads))
	}
	if model.err != nil {
		t.Errorf("expected nil error, got %v", model.err)
	}
}

func TestThreadNavigation(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewThreads
	m.loading = false
	m.threads = []sdk.Thread{
		{ThreadID: "t-1", Title: "One"},
		{ThreadID: "t-2", Title: "Two"},
	}
	m.threadCursor = 0 // "New Thread" position

	// Move down to first thread
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := updated.(Model)
	if model.threadCursor != 1 {
		t.Errorf("after j: expected cursor=1, got %d", model.threadCursor)
	}

	// Move down to second thread
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(Model)
	if model.threadCursor != 2 {
		t.Errorf("after j: expected cursor=2, got %d", model.threadCursor)
	}

	// Bounds check: can't go past last
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(Model)
	if model.threadCursor != 2 {
		t.Errorf("after j at end: expected cursor=2, got %d", model.threadCursor)
	}

	// Move back up
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.threadCursor != 1 {
		t.Errorf("after k: expected cursor=1, got %d", model.threadCursor)
	}
}

func TestThreadsView_Error(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewThreads
	m.loading = true

	updated, _ := m.Update(errMsg{err: errors.New("thread load failed")})
	model := updated.(Model)

	if model.loading {
		t.Error("expected loading=false after errMsg")
	}

	output := model.View()
	if !strings.Contains(output, "thread load failed") {
		t.Errorf("expected View output to contain error message, got: %s", output)
	}
}

func TestNewThreadSelection(t *testing.T) {
	cfg := &config.Config{PAT: "test-pat"}
	m := New(cfg)
	m.view = viewThreads
	m.loading = false
	m.threads = []sdk.Thread{{ThreadID: "t-1", Title: "Existing"}}
	m.threadCursor = 0 // "New Thread" is at position 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)

	if model.view != viewChat {
		t.Errorf("expected viewChat, got %d", model.view)
	}
	if model.threadID != "" {
		t.Errorf("expected empty threadID for new thread, got %q", model.threadID)
	}
	if model.activeThread != nil {
		t.Error("expected nil activeThread for new thread")
	}
}
