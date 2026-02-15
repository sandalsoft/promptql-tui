package sdk

import (
	"errors"
	"net/http"
	"testing"
)

// ---------------------------------------------------------------------------
// ListUserProjects
// ---------------------------------------------------------------------------

func TestListUserProjects_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		body := graphqlJSON(`{
			"ddn_projects": [
				{"id": "proj-1", "name": "Project Alpha", "ddn_builds": [{"fqdn": "alpha.ddn.hasura.app"}]},
				{"id": "proj-2", "name": "Project Beta", "ddn_builds": []}
			]
		}`)
		return jsonResponse(200, body), nil
	})

	projects, err := client.Projects().ListUserProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	// First project — has a build FQDN
	if projects[0].Name != "Project Alpha" {
		t.Errorf("expected name 'Project Alpha', got %q", projects[0].Name)
	}
	if projects[0].DDNProjectID != "proj-1" {
		t.Errorf("expected DDNProjectID 'proj-1', got %q", projects[0].DDNProjectID)
	}
	if projects[0].BuildFQDN != "alpha.ddn.hasura.app" {
		t.Errorf("expected BuildFQDN 'alpha.ddn.hasura.app', got %q", projects[0].BuildFQDN)
	}

	// Second project — no builds
	if projects[1].Name != "Project Beta" {
		t.Errorf("expected name 'Project Beta', got %q", projects[1].Name)
	}
	if projects[1].BuildFQDN != "" {
		t.Errorf("expected empty BuildFQDN, got %q", projects[1].BuildFQDN)
	}
}

func TestListUserProjects_AuthError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(401, `{"message":"invalid token"}`), nil
	})

	_, err := client.Projects().ListUserProjects()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Errorf("expected *AuthenticationError, got %T: %v", err, err)
	}
}

func TestListUserProjects_RuntimeError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	_, err := client.Projects().ListUserProjects()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "connection refused") {
		t.Errorf("expected error to contain 'connection refused', got: %v", err)
	}
}

func TestListUserProjects_MalformedResponse(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{not valid json`), nil
	})

	_, err := client.Projects().ListUserProjects()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---------------------------------------------------------------------------
// GetConfig
// ---------------------------------------------------------------------------

func TestGetConfig_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		body := graphqlJSON(`{
			"getPromptQlConfig": {
				"promptQlEnabled": true,
				"playgroundEnabled": false
			}
		}`)
		return jsonResponse(200, body), nil
	})

	cfg, err := client.Projects().GetConfig("proj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.PromptQLEnabled {
		t.Error("expected PromptQLEnabled to be true")
	}
	if cfg.PlaygroundEnabled {
		t.Error("expected PlaygroundEnabled to be false")
	}
}

func TestGetConfig_AuthError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(401, `{"message":"unauthorized"}`), nil
	})

	_, err := client.Projects().GetConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Errorf("expected *AuthenticationError, got %T: %v", err, err)
	}
}

func TestGetConfig_RuntimeError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network timeout")
	})

	_, err := client.Projects().GetConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "network timeout") {
		t.Errorf("expected error to contain 'network timeout', got: %v", err)
	}
}

func TestGetConfig_MalformedResponse(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `<<<broken>>>`), nil
	})

	_, err := client.Projects().GetConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---------------------------------------------------------------------------
// GetPlaygroundConfig
// ---------------------------------------------------------------------------

func TestGetPlaygroundConfig_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		body := graphqlJSON(`{
			"getPlaygroundConfig": {
				"allowPublicAccess": true,
				"featureFlags": {"beta": true},
				"llmApiKey": "sk-test-key",
				"llmProvider": "openai",
				"projectTokenUsageLimit": 50000,
				"readme": "# Welcome",
				"systemInstructions": "Be helpful",
				"userTokenUsageLimit": 10000
			}
		}`)
		return jsonResponse(200, body), nil
	})

	cfg, err := client.Projects().GetPlaygroundConfig("proj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AllowPublicAccess == nil || !*cfg.AllowPublicAccess {
		t.Error("expected AllowPublicAccess to be true")
	}
	if cfg.LLMApiKey != "sk-test-key" {
		t.Errorf("expected LLMApiKey 'sk-test-key', got %q", cfg.LLMApiKey)
	}
	if cfg.LLMProvider != "openai" {
		t.Errorf("expected LLMProvider 'openai', got %q", cfg.LLMProvider)
	}
	if cfg.ProjectTokenUsageLimit == nil || *cfg.ProjectTokenUsageLimit != 50000 {
		t.Errorf("expected ProjectTokenUsageLimit 50000, got %v", cfg.ProjectTokenUsageLimit)
	}
	if cfg.Readme != "# Welcome" {
		t.Errorf("expected Readme '# Welcome', got %q", cfg.Readme)
	}
	if cfg.SystemInstructions != "Be helpful" {
		t.Errorf("expected SystemInstructions 'Be helpful', got %q", cfg.SystemInstructions)
	}
	if cfg.UserTokenUsageLimit == nil || *cfg.UserTokenUsageLimit != 10000 {
		t.Errorf("expected UserTokenUsageLimit 10000, got %v", cfg.UserTokenUsageLimit)
	}
	if cfg.FeatureFlags == nil {
		t.Fatal("expected FeatureFlags to be non-nil")
	}
	if beta, ok := cfg.FeatureFlags["beta"].(bool); !ok || !beta {
		t.Errorf("expected FeatureFlags['beta'] to be true, got %v", cfg.FeatureFlags["beta"])
	}
}

func TestGetPlaygroundConfig_AuthError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(401, `{"message":"unauthorized"}`), nil
	})

	_, err := client.Projects().GetPlaygroundConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Errorf("expected *AuthenticationError, got %T: %v", err, err)
	}
}

func TestGetPlaygroundConfig_RuntimeError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("dns resolution failed")
	})

	_, err := client.Projects().GetPlaygroundConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "dns resolution failed") {
		t.Errorf("expected error to contain 'dns resolution failed', got: %v", err)
	}
}

func TestGetPlaygroundConfig_MalformedResponse(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `not json at all`), nil
	})

	_, err := client.Projects().GetPlaygroundConfig("proj-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
