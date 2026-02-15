// Package sdk provides a Go client for the PromptQL API.
// Vendored from github.com/sandalsoft/promptql-sdk-golang
package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientOptions configures the PromptQL client.
type ClientOptions struct {
	// PAT is the Personal Access Token for control-plane operations.
	PAT string
	// APIKey is the API key for the natural-language query endpoint.
	APIKey string
	// ProjectID is the default project ID used when not specified per-call.
	ProjectID string
	// BaseURL is the base URL of the GraphQL control-plane API.
	BaseURL string
	// APIURL is the base URL of the REST natural-language query API.
	APIURL string
	// AuthURL is the base URL of the authentication service.
	AuthURL string
	// ControlPlaneURL is the base URL of the DDN control-plane API for project listing.
	ControlPlaneURL string
	// Timeout is the HTTP request timeout.
	Timeout time.Duration
	// HTTPClient allows injecting a custom *http.Client (useful for testing).
	HTTPClient *http.Client
}

// Client is the main entry point for the PromptQL SDK.
type Client struct {
	pat             string
	apiKey          string
	ProjectID       string
	baseURL         string
	apiURL          string
	authURL         string
	controlPlaneURL string
	http            *http.Client

	projects *ProjectsResource
	prompts  *PromptsResource
	apiKeys  *APIKeysResource
	threads  *ThreadsResource
	query    *QueryResource
	users    *UsersResource
}

// NewClient creates a new PromptQL client with the given options.
func NewClient(opts ClientOptions) *Client {
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = "https://data.promptql.pro.hasura.io"
	}
	apiURL := opts.APIURL
	if apiURL == "" {
		apiURL = "https://api.promptql.pro.hasura.io"
	}
	authURL := opts.AuthURL
	if authURL == "" {
		authURL = "https://auth.pro.hasura.io"
	}
	controlPlaneURL := opts.ControlPlaneURL
	if controlPlaneURL == "" {
		controlPlaneURL = "https://data.pro.hasura.io"
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}

	c := &Client{
		pat:             opts.PAT,
		apiKey:          opts.APIKey,
		ProjectID:       opts.ProjectID,
		baseURL:         baseURL,
		apiURL:          apiURL,
		authURL:         authURL,
		controlPlaneURL: controlPlaneURL,
		http:            httpClient,
	}

	c.projects = &ProjectsResource{client: c}
	c.prompts = &PromptsResource{client: c}
	c.apiKeys = &APIKeysResource{client: c}
	c.threads = &ThreadsResource{client: c}
	c.query = &QueryResource{client: c}
	c.users = &UsersResource{client: c}

	return c
}

// Projects returns the projects resource.
func (c *Client) Projects() *ProjectsResource { return c.projects }

// Prompts returns the prompts resource.
func (c *Client) Prompts() *PromptsResource { return c.prompts }

// APIKeys returns the API keys resource.
func (c *Client) APIKeys() *APIKeysResource { return c.apiKeys }

// Threads returns the threads resource.
func (c *Client) Threads() *ThreadsResource { return c.threads }

// Query returns the query resource.
func (c *Client) Query() *QueryResource { return c.query }

// Users returns the users resource.
func (c *Client) Users() *UsersResource { return c.users }

// GetDDNToken exchanges a PAT for a DDN bearer token.
func (c *Client) GetDDNToken(projectID string) (*TokenResponse, error) {
	if c.pat == "" {
		return nil, &AuthenticationError{PromptQLError{
			Message: "A Personal Access Token (pat) is required to obtain a DDN token",
		}}
	}

	req, err := http.NewRequest("POST", c.authURL+"/ddn/promptql/token", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "pat "+c.pat)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-project-id", projectID)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &token, nil
}

// authHeader returns the Authorization header value for the given auth type.
func (c *Client) authHeader(authType string) (string, error) {
	switch authType {
	case "pat":
		if c.pat == "" {
			return "", &AuthenticationError{PromptQLError{
				Message: "A Personal Access Token (pat) is required for this operation",
			}}
		}
		return "pat " + c.pat, nil
	case "bearer":
		if c.apiKey == "" {
			return "", &AuthenticationError{PromptQLError{
				Message: "An API key is required for this operation",
			}}
		}
		return "Bearer " + c.apiKey, nil
	default:
		return "", fmt.Errorf("unknown auth type: %s", authType)
	}
}

// graphqlRequest is the payload sent for GraphQL operations.
type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// graphqlResponse is the raw GraphQL response.
type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// GraphQL executes a GraphQL query/mutation and returns the data field.
func (c *Client) GraphQL(query string, variables map[string]interface{}, authType string) (map[string]interface{}, error) {
	return c.graphqlTo(c.baseURL+"/graphql", query, variables, authType)
}

// GraphQLControlPlane executes a GraphQL query against the DDN control-plane API.
func (c *Client) GraphQLControlPlane(query string, variables map[string]interface{}, authType string) (map[string]interface{}, error) {
	return c.graphqlTo(c.controlPlaneURL+"/v1/graphql", query, variables, authType)
}

// graphqlTo executes a GraphQL query/mutation against the given endpoint URL.
func (c *Client) graphqlTo(endpoint string, query string, variables map[string]interface{}, authType string) (map[string]interface{}, error) {
	auth, err := c.authHeader(authType)
	if err != nil {
		return nil, err
	}

	payload := graphqlRequest{Query: query, Variables: variables}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var gqlResp graphqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, &PromptQLError{Message: gqlResp.Errors[0].Message}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(gqlResp.Data, &data); err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return data, nil
}

// PostAPI sends a POST request to the REST API and returns the response body.
func (c *Client) PostAPI(path string, jsonBody map[string]interface{}, authType string) (map[string]interface{}, error) {
	auth, err := c.authHeader(authType)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}

// checkResponse inspects the HTTP response and returns a typed error for non-2xx status codes.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var parsed map[string]interface{}
	message := string(body)
	if err := json.Unmarshal(body, &parsed); err == nil {
		if msg, ok := parsed["message"].(string); ok {
			message = msg
		} else if msg, ok := parsed["error"].(string); ok {
			message = msg
		}
	}

	return newErrorForStatus(resp.StatusCode, message)
}

// decodeJSONField is a helper that extracts a named field from a map and decodes it into the target.
func decodeJSONField(data map[string]interface{}, field string, target interface{}) error {
	val, ok := data[field]
	if !ok {
		return fmt.Errorf("field %q not found in response", field)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshaling field %q: %w", field, err)
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("unmarshaling field %q: %w", field, err)
	}
	return nil
}
