package sdk

// APIKeysResource provides management of runtime API keys.
type APIKeysResource struct {
	client *Client
}

// List lists all runtime API keys for a project.
func (r *APIKeysResource) List(projectID string) ([]RuntimeAPIKey, error) {
	query := `
		query ListRuntimeApiKeys($projectId: String!) {
			getRuntimeApiKeys(projectId: $projectId) {
				id
				name
				projectId
				apiKeyMasked
				isActive
				createdAt
				createdBy
				lastUsedAt
				promptqlTimeout
				sqlTimeout
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result []RuntimeAPIKey
	if err := decodeJSONField(data, "getRuntimeApiKeys", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateOptions configures API key generation.
type GenerateOptions struct {
	ProjectID       string
	Name            string
	PromptQLTimeout *int
	SQLTimeout      *int
}

// Generate creates a new runtime API key.
func (r *APIKeysResource) Generate(opts GenerateOptions) (map[string]interface{}, error) {
	query := `
		mutation GenerateRuntimeApiKey($projectId: String!, $name: String!, $promptqlTimeout: Int, $sqlTimeout: Int) {
			generateRuntimeApiKey(projectId: $projectId, name: $name, promptqlTimeout: $promptqlTimeout, sqlTimeout: $sqlTimeout) {
				id
				name
				projectId
				apiKey
				apiKeyMasked
				isActive
				createdAt
				createdBy
				promptqlTimeout
				sqlTimeout
			}
		}`
	variables := map[string]interface{}{
		"projectId": opts.ProjectID,
		"name":      opts.Name,
	}
	if opts.PromptQLTimeout != nil {
		variables["promptqlTimeout"] = *opts.PromptQLTimeout
	}
	if opts.SQLTimeout != nil {
		variables["sqlTimeout"] = *opts.SQLTimeout
	}

	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	result, ok := data["generateRuntimeApiKey"].(map[string]interface{})
	if !ok {
		return nil, &PromptQLError{Message: "unexpected response format for generateRuntimeApiKey"}
	}
	return result, nil
}

// Remove removes (deactivates) a runtime API key.
func (r *APIKeysResource) Remove(projectID string, apiKeyID int) (*MessageResult, error) {
	query := `
		mutation RemoveRuntimeApiKey($projectId: String!, $apiKeyId: Int!) {
			removeRuntimeApiKey(projectId: $projectId, apiKeyId: $apiKeyId) {
				message
			}
		}`
	variables := map[string]interface{}{
		"projectId": projectID,
		"apiKeyId":  apiKeyID,
	}
	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result MessageResult
	if err := decodeJSONField(data, "removeRuntimeApiKey", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
