package sdk

// PromptsResource provides CRUD operations on sample prompts.
type PromptsResource struct {
	client *Client
}

// List lists all sample prompts for a project.
func (r *PromptsResource) List(projectID string) ([]SamplePrompt, error) {
	query := `
		query ListSamplePrompts($projectId: String!) {
			getSamplePrompts(projectId: $projectId) {
				id
				displayText
				fullPrompt
				projectId
				createdBy
				createdAt
				updatedBy
				updatedAt
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result []SamplePrompt
	if err := decodeJSONField(data, "getSamplePrompts", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create creates a new sample prompt.
func (r *PromptsResource) Create(projectID, displayText, fullPrompt string) (*SamplePrompt, error) {
	query := `
		mutation CreateSamplePrompt($projectId: String!, $displayText: String!, $fullPrompt: String!) {
			createSamplePrompt(projectId: $projectId, displayText: $displayText, fullPrompt: $fullPrompt) {
				id
				displayText
				fullPrompt
				projectId
				createdBy
				createdAt
				updatedBy
				updatedAt
			}
		}`
	variables := map[string]interface{}{
		"projectId":   projectID,
		"displayText": displayText,
		"fullPrompt":  fullPrompt,
	}
	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result SamplePrompt
	if err := decodeJSONField(data, "createSamplePrompt", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing sample prompt.
func (r *PromptsResource) Update(projectID, promptID, displayText, fullPrompt string) (*SamplePrompt, error) {
	query := `
		mutation UpdateSamplePrompt($projectId: String!, $promptId: String!, $displayText: String!, $fullPrompt: String!) {
			updateSamplePrompt(projectId: $projectId, promptId: $promptId, displayText: $displayText, fullPrompt: $fullPrompt) {
				id
				displayText
				fullPrompt
				projectId
				createdBy
				createdAt
				updatedBy
				updatedAt
			}
		}`
	variables := map[string]interface{}{
		"projectId":   projectID,
		"promptId":    promptID,
		"displayText": displayText,
		"fullPrompt":  fullPrompt,
	}
	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result SamplePrompt
	if err := decodeJSONField(data, "updateSamplePrompt", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a sample prompt.
func (r *PromptsResource) Delete(projectID, promptID string) (*MessageResult, error) {
	query := `
		mutation DeleteSamplePrompt($projectId: String!, $promptId: String!) {
			deleteSamplePrompt(projectId: $projectId, promptId: $promptId) {
				message
			}
		}`
	variables := map[string]interface{}{
		"projectId": projectID,
		"promptId":  promptID,
	}
	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result MessageResult
	if err := decodeJSONField(data, "deleteSamplePrompt", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
