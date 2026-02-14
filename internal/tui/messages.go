package tui

import "github.com/sandalsoft/promptql-tui/internal/sdk"

// Message types for the TUI event loop.

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type projectsLoadedMsg struct {
	projects []sdk.UserProject
}

type projectSelectedMsg struct {
	project sdk.UserProject
}

type threadsLoadedMsg struct {
	threads []sdk.Thread
}

type threadStartedMsg struct {
	result *sdk.StartThreadResult
}

type messageSentMsg struct {
	result *sdk.SendMessageResult
}

type eventsLoadedMsg struct {
	events []sdk.ThreadEvent
}

type queryResultMsg struct {
	result map[string]interface{}
}

type configSavedMsg struct{}

type lookupResultMsg struct {
	result *sdk.LookupProjectResult
}
