package sdk

// QueryResource provides natural language query execution.
type QueryResource struct {
	client *Client
}

// ExecuteOptions configures a natural language query.
type ExecuteOptions struct {
	Interactions []map[string]interface{}
	DDNURL       string
	DDNHeaders   map[string]string
	Timezone     string // defaults to "UTC"
	Stream       bool
	Version      string // defaults to "v1"
}

// Execute sends a natural language query to the PromptQL query endpoint.
func (r *QueryResource) Execute(opts ExecuteOptions) (map[string]interface{}, error) {
	tz := opts.Timezone
	if tz == "" {
		tz = "UTC"
	}
	version := opts.Version
	if version == "" {
		version = "v1"
	}

	ddn := map[string]interface{}{"url": opts.DDNURL}
	if opts.DDNHeaders != nil {
		headers := map[string]interface{}{}
		for k, v := range opts.DDNHeaders {
			headers[k] = v
		}
		ddn["headers"] = headers
	}

	body := map[string]interface{}{
		"version":      version,
		"stream":       opts.Stream,
		"timezone":     tz,
		"interactions": opts.Interactions,
		"ddn":          ddn,
	}

	return r.client.PostAPI("/query", body, "bearer")
}

// Ask is a convenience method that sends a single user question.
func (r *QueryResource) Ask(question, ddnURL string, ddnHeaders map[string]string, timezone string) (map[string]interface{}, error) {
	if timezone == "" {
		timezone = "UTC"
	}
	interactions := []map[string]interface{}{
		{
			"role": "user",
			"user_message": map[string]interface{}{
				"text": question,
			},
		},
	}
	return r.Execute(ExecuteOptions{
		Interactions: interactions,
		DDNURL:       ddnURL,
		DDNHeaders:   ddnHeaders,
		Timezone:     timezone,
	})
}
