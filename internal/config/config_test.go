package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoFile(t *testing.T) {
	// Point HOME at an empty temp dir so configPath resolves to a non-existent file.
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PAT != "" {
		t.Errorf("expected empty PAT, got %q", cfg.PAT)
	}
	if cfg.APIKey != "" {
		t.Errorf("expected empty APIKey, got %q", cfg.APIKey)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "promptql-tui")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("creating config dir: %v", err)
	}

	want := Config{
		PAT:       "my-pat",
		APIKey:    "my-key",
		ProjectID: "proj-123",
		DDNURL:    "https://example.ddn.hasura.app/graphql",
		Timezone:  "America/New_York",
	}
	data, _ := json.MarshalIndent(want, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PAT != want.PAT {
		t.Errorf("PAT: got %q, want %q", cfg.PAT, want.PAT)
	}
	if cfg.APIKey != want.APIKey {
		t.Errorf("APIKey: got %q, want %q", cfg.APIKey, want.APIKey)
	}
	if cfg.ProjectID != want.ProjectID {
		t.Errorf("ProjectID: got %q, want %q", cfg.ProjectID, want.ProjectID)
	}
	if cfg.DDNURL != want.DDNURL {
		t.Errorf("DDNURL: got %q, want %q", cfg.DDNURL, want.DDNURL)
	}
	if cfg.Timezone != want.Timezone {
		t.Errorf("Timezone: got %q, want %q", cfg.Timezone, want.Timezone)
	}
}

func TestLoad_MalformedJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "promptql-tui")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("creating config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{not json}`), 0600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestHasCredentials_WithPAT(t *testing.T) {
	cfg := Config{PAT: "some-pat"}
	if !cfg.HasCredentials() {
		t.Error("expected HasCredentials() to return true with PAT set")
	}
}

func TestHasCredentials_WithoutPAT(t *testing.T) {
	cfg := Config{}
	if cfg.HasCredentials() {
		t.Error("expected HasCredentials() to return false with empty PAT")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	original := &Config{
		PAT:       "round-trip-pat",
		APIKey:    "round-trip-key",
		ProjectID: "proj-rt",
		DDNURL:    "https://rt.example.com/graphql",
		Timezone:  "Europe/London",
	}

	if err := original.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.PAT != original.PAT {
		t.Errorf("PAT: got %q, want %q", loaded.PAT, original.PAT)
	}
	if loaded.APIKey != original.APIKey {
		t.Errorf("APIKey: got %q, want %q", loaded.APIKey, original.APIKey)
	}
	if loaded.ProjectID != original.ProjectID {
		t.Errorf("ProjectID: got %q, want %q", loaded.ProjectID, original.ProjectID)
	}
	if loaded.DDNURL != original.DDNURL {
		t.Errorf("DDNURL: got %q, want %q", loaded.DDNURL, original.DDNURL)
	}
	if loaded.Timezone != original.Timezone {
		t.Errorf("Timezone: got %q, want %q", loaded.Timezone, original.Timezone)
	}
}
