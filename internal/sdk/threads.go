package sdk

// ThreadsResource provides thread and conversation management.
type ThreadsResource struct {
	client *Client
}

// StartOptions configures thread creation.
type StartOptions struct {
	ProjectID  string
	Message    string
	BuildFQDN  string
	Timezone   string // defaults to "UTC"
	Visibility string // optional
}

// Start starts a new thread with an initial message.
func (r *ThreadsResource) Start(opts StartOptions) (*StartThreadResult, error) {
	query := `
		mutation StartThread($projectId: String!, $message: String!, $buildFqdn: String!, $timezone: String!, $visibility: String) {
			startThread(projectId: $projectId, message: $message, buildFqdn: $buildFqdn, timezone: $timezone, visibility: $visibility) {
				thread_id
				title
				created_at
				updated_at
				thread_events {
					thread_event_id
					thread_id
					event_data
					created_at
					user_id
				}
			}
		}`
	tz := opts.Timezone
	if tz == "" {
		tz = "UTC"
	}
	variables := map[string]interface{}{
		"projectId": opts.ProjectID,
		"message":   opts.Message,
		"buildFqdn": opts.BuildFQDN,
		"timezone":  tz,
	}
	if opts.Visibility != "" {
		variables["visibility"] = opts.Visibility
	}

	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result StartThreadResult
	if err := decodeJSONField(data, "startThread", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendMessageOptions configures sending a message in a thread.
type SendMessageOptions struct {
	ThreadID  string
	Message   string
	BuildFQDN string
	Timezone  string // defaults to "UTC"
}

// SendMessage sends a follow-up message to an existing thread.
func (r *ThreadsResource) SendMessage(opts SendMessageOptions) (*SendMessageResult, error) {
	query := `
		mutation SendMessage($threadId: String!, $message: String!, $buildFqdn: String!, $timezone: String!) {
			sendMessage(threadId: $threadId, message: $message, buildFqdn: $buildFqdn, timezone: $timezone) {
				thread_event_id
				event_data
				created_at
			}
		}`
	tz := opts.Timezone
	if tz == "" {
		tz = "UTC"
	}
	variables := map[string]interface{}{
		"threadId":  opts.ThreadID,
		"message":   opts.Message,
		"buildFqdn": opts.BuildFQDN,
		"timezone":  tz,
	}

	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result SendMessageResult
	if err := decodeJSONField(data, "sendMessage", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get fetches a single thread by ID.
func (r *ThreadsResource) Get(threadID string) (*Thread, error) {
	query := `
		query GetThread($threadId: String!) {
			getThread(threadId: $threadId) {
				thread_id
				title
				created_at
				updated_at
				project_id
				build_id
				user_id
				visibility
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"threadId": threadID}, "pat")
	if err != nil {
		return nil, err
	}
	var result Thread
	if err := decodeJSONField(data, "getThread", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List lists threads for a project and user.
func (r *ThreadsResource) List(projectID, userID string) ([]Thread, error) {
	query := `
		query ListThreads($projectId: String!, $userId: String!) {
			getThreads(projectId: $projectId, userId: $userId) {
				thread_id
				title
				created_at
				updated_at
				project_id
				build_id
				user_id
				visibility
			}
		}`
	variables := map[string]interface{}{
		"projectId": projectID,
		"userId":    userID,
	}
	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result []Thread
	if err := decodeJSONField(data, "getThreads", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetEvents fetches all events for a thread.
func (r *ThreadsResource) GetEvents(threadID string) ([]ThreadEvent, error) {
	query := `
		query GetThreadEvents($threadId: String!) {
			getThreadEvents(threadId: $threadId) {
				thread_event_id
				thread_id
				event_data
				created_at
				user_id
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"threadId": threadID}, "pat")
	if err != nil {
		return nil, err
	}
	var result []ThreadEvent
	if err := decodeJSONField(data, "getThreadEvents", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SubmitFeedback submits feedback for a thread message.
func (r *ThreadsResource) SubmitFeedback(threadID, messageID string, feedback int, details string) (*ThreadFeedback, error) {
	query := `
		mutation SubmitFeedback($threadId: String!, $messageId: String!, $feedback: Int!, $details: String) {
			submitThreadFeedback(threadId: $threadId, messageId: $messageId, feedback: $feedback, details: $details) {
				thread_id
				message_id
				promptql_user_id
				feedback
				details
				created_at
			}
		}`
	variables := map[string]interface{}{
		"threadId":  threadID,
		"messageId": messageID,
		"feedback":  feedback,
	}
	if details != "" {
		variables["details"] = details
	}

	data, err := r.client.GraphQL(query, variables, "pat")
	if err != nil {
		return nil, err
	}
	var result ThreadFeedback
	if err := decodeJSONField(data, "submitThreadFeedback", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
