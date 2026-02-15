package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sandalsoft/promptql-tui/internal/config"
	"github.com/sandalsoft/promptql-tui/internal/tui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Allow overriding PAT from environment
	if pat := os.Getenv("PROMPTQL_PAT"); pat != "" {
		cfg.PAT = pat
	}
	if apiKey := os.Getenv("PROMPTQL_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if ddnURL := os.Getenv("PROMPTQL_DDN_URL"); ddnURL != "" {
		cfg.DDNURL = ddnURL
	}

	m := tui.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
