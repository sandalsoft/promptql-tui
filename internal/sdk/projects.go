package sdk

// ProjectsResource provides operations on PromptQL projects.
type ProjectsResource struct {
	client *Client
}

// GetConfig fetches the PromptQL feature configuration for a project.
func (r *ProjectsResource) GetConfig(projectID string) (*PromptQLConfig, error) {
	query := `
		query GetPromptQLConfig($projectId: String!) {
			getPromptQlConfig(projectId: $projectId) {
				promptQlEnabled
				playgroundEnabled
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result PromptQLConfig
	if err := decodeJSONField(data, "getPromptQlConfig", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlaygroundConfig fetches the playground configuration for a project.
func (r *ProjectsResource) GetPlaygroundConfig(projectID string) (*PlaygroundConfig, error) {
	query := `
		query GetPlaygroundConfig($projectId: String!) {
			getPlaygroundConfig(projectId: $projectId) {
				allowPublicAccess
				featureFlags
				llmApiKey
				llmProvider
				projectTokenUsageLimit
				readme
				systemInstructions
				userTokenUsageLimit
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result PlaygroundConfig
	if err := decodeJSONField(data, "getPlaygroundConfig", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserProjects lists all projects visible to the authenticated user.
func (r *ProjectsResource) ListUserProjects() ([]UserProject, error) {
	query := `
		query ListUserProjects {
			getUserProjects {
				buildFqdn
				ddnProjectId
				name
				projectId
			}
		}`
	data, err := r.client.GraphQL(query, nil, "pat")
	if err != nil {
		return nil, err
	}
	var result []UserProject
	if err := decodeJSONField(data, "getUserProjects", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// LookupOptions configures a project lookup.
type LookupOptions struct {
	ProjectID   string
	ProjectName string
	FQDN        string
}

// Lookup looks up a project by ID, name, or FQDN.
func (r *ProjectsResource) Lookup(opts LookupOptions) (*LookupProjectResult, error) {
	query := `
		query LookupProject($projectId: String, $projectName: String, $fqdn: String) {
			lookupProject(projectId: $projectId, projectName: $projectName, fqdn: $fqdn) {
				buildFqdn
				consoleUrl
				ddnProjectId
				name
				projectId
			}
		}`
	variables := map[string]interface{}{}
	if opts.ProjectID != "" {
		variables["projectId"] = opts.ProjectID
	}
	if opts.ProjectName != "" {
		variables["projectName"] = opts.ProjectName
	}
	if opts.FQDN != "" {
		variables["fqdn"] = opts.FQDN
	}

	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result LookupProjectResult
	if err := decodeJSONField(data, "lookupProject", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Enable enables PromptQL for a project.
func (r *ProjectsResource) Enable(projectID string) (*MessageResult, error) {
	query := `
		mutation EnablePromptQL($projectId: String!) {
			enablePromptQl(projectId: $projectId) {
				message
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result MessageResult
	if err := decodeJSONField(data, "enablePromptQl", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Disable disables PromptQL for a project.
func (r *ProjectsResource) Disable(projectID string) (*MessageResult, error) {
	query := `
		mutation DisablePromptQL($projectId: String!) {
			disablePromptQl(projectId: $projectId) {
				message
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"projectId": projectID}, "pat")
	if err != nil {
		return nil, err
	}
	var result MessageResult
	if err := decodeJSONField(data, "disablePromptQl", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
