package sdk

// TokenResponse is returned from the DDN token exchange endpoint.
type TokenResponse struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
	Status string `json:"status,omitempty"`
}

// PromptQLConfig holds the PromptQL feature configuration for a project.
type PromptQLConfig struct {
	PromptQLEnabled   bool `json:"promptQlEnabled"`
	PlaygroundEnabled bool `json:"playgroundEnabled"`
}

// SamplePrompt is a sample prompt associated with a project.
type SamplePrompt struct {
	ID          string `json:"id"`
	DisplayText string `json:"displayText"`
	FullPrompt  string `json:"fullPrompt"`
	ProjectID   string `json:"projectId"`
	CreatedBy   string `json:"createdBy,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedBy   string `json:"updatedBy,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// RuntimeAPIKey is a runtime API key for the query endpoint.
type RuntimeAPIKey struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ProjectID       string `json:"projectId"`
	APIKeyMasked    string `json:"apiKeyMasked,omitempty"`
	IsActive        *bool  `json:"isActive,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	CreatedBy       string `json:"createdBy,omitempty"`
	LastUsedAt      string `json:"lastUsedAt,omitempty"`
	PromptQLTimeout *int   `json:"promptqlTimeout,omitempty"`
	SQLTimeout      *int   `json:"sqlTimeout,omitempty"`
}

// UserProject is a project visible to the current user.
type UserProject struct {
	BuildFQDN    string `json:"buildFqdn,omitempty"`
	DDNProjectID string `json:"ddnProjectId,omitempty"`
	Name         string `json:"name"`
	ProjectID    string `json:"projectId"`
}

// PlaygroundConfig holds playground configuration for a project.
type PlaygroundConfig struct {
	AllowPublicAccess      *bool                  `json:"allowPublicAccess,omitempty"`
	FeatureFlags           map[string]interface{} `json:"featureFlags,omitempty"`
	LLMApiKey              string                 `json:"llmApiKey,omitempty"`
	LLMProvider            string                 `json:"llmProvider,omitempty"`
	ProjectTokenUsageLimit *int                   `json:"projectTokenUsageLimit,omitempty"`
	Readme                 string                 `json:"readme,omitempty"`
	SystemInstructions     string                 `json:"systemInstructions,omitempty"`
	UserTokenUsageLimit    *int                   `json:"userTokenUsageLimit,omitempty"`
}

// LookupProjectResult is the result of looking up a project.
type LookupProjectResult struct {
	BuildFQDN    string `json:"buildFqdn,omitempty"`
	ConsoleURL   string `json:"consoleUrl,omitempty"`
	DDNProjectID string `json:"ddnProjectId,omitempty"`
	Name         string `json:"name"`
	ProjectID    string `json:"projectId"`
}

// ThreadEvent is an event within a thread.
type ThreadEvent struct {
	ThreadEventID int                    `json:"thread_event_id"`
	ThreadID      string                 `json:"thread_id,omitempty"`
	EventData     map[string]interface{} `json:"event_data,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
}

// Thread is a conversation thread.
type Thread struct {
	ThreadID   string `json:"thread_id"`
	Title      string `json:"title,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	ProjectID  string `json:"project_id,omitempty"`
	BuildID    string `json:"build_id,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	Visibility string `json:"visibility,omitempty"`
}

// StartThreadResult is the result of starting a new thread.
type StartThreadResult struct {
	ThreadID     string        `json:"thread_id"`
	Title        string        `json:"title,omitempty"`
	CreatedAt    string        `json:"created_at,omitempty"`
	UpdatedAt    string        `json:"updated_at,omitempty"`
	ThreadEvents []ThreadEvent `json:"thread_events,omitempty"`
}

// SendMessageResult is the result of sending a message to a thread.
type SendMessageResult struct {
	ThreadEventID int                    `json:"thread_event_id"`
	EventData     map[string]interface{} `json:"event_data,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
}

// ThreadFeedback holds feedback on a thread message.
type ThreadFeedback struct {
	ThreadID       string `json:"thread_id"`
	MessageID      string `json:"message_id"`
	PromptQLUserID string `json:"promptql_user_id,omitempty"`
	Feedback       *int   `json:"feedback,omitempty"`
	Details        string `json:"details,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

// PromptQLUser represents a PromptQL user.
type PromptQLUser struct {
	PromptQLUserID     string `json:"promptql_user_id"`
	ControlPlaneUserID string `json:"control_plane_user_id,omitempty"`
	Email              string `json:"email,omitempty"`
	DisplayName        string `json:"display_name,omitempty"`
	IsActive           *bool  `json:"is_active,omitempty"`
	ProjectID          string `json:"project_id,omitempty"`
}

// Program represents a PromptQL program.
type Program struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// MessageResult is a generic message result from mutation operations.
type MessageResult struct {
	Message string `json:"message"`
}
