package sdk

// UsersResource provides operations on PromptQL user accounts.
type UsersResource struct {
	client *Client
}

// GetCurrent fetches a PromptQL user by their control-plane user ID.
func (r *UsersResource) GetCurrent(controlPlaneUserID string) (*PromptQLUser, error) {
	query := `
		query GetPromptQLUser($controlPlaneUserId: String!) {
			getPromptQLUser(controlPlaneUserId: $controlPlaneUserId) {
				promptql_user_id
				control_plane_user_id
				email
				display_name
				is_active
				project_id
			}
		}`
	data, err := r.client.GraphQL(query, map[string]interface{}{"controlPlaneUserId": controlPlaneUserID}, "pat")
	if err != nil {
		return nil, err
	}
	var result PromptQLUser
	if err := decodeJSONField(data, "getPromptQLUser", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List lists all PromptQL users.
func (r *UsersResource) List() ([]PromptQLUser, error) {
	query := `
		query ListPromptQLUsers {
			getPromptQLUsers {
				promptql_user_id
				control_plane_user_id
				email
				display_name
				is_active
				project_id
			}
		}`
	data, err := r.client.GraphQL(query, nil, "pat")
	if err != nil {
		return nil, err
	}
	var result []PromptQLUser
	if err := decodeJSONField(data, "getPromptQLUsers", &result); err != nil {
		return nil, err
	}
	return result, nil
}
