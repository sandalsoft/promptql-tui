package sdk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// mockRoundTripper implements http.RoundTripper for testing.
// It delegates to a configurable function so each test can define
// the exact HTTP response (or error) to return.
type mockRoundTripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

// newTestClient creates an sdk.Client with the given mock transport function.
// The PAT is pre-set so auth-required methods work by default.
func newTestClient(fn func(*http.Request) (*http.Response, error)) *Client {
	return NewClient(ClientOptions{
		PAT:             "test-pat",
		BaseURL:         "https://test.example.com",
		APIURL:          "https://api.test.example.com",
		AuthURL:         "https://auth.test.example.com",
		ControlPlaneURL: "https://cp.test.example.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{fn: fn},
		},
	})
}

// jsonResponse builds a *http.Response with the given status code and JSON body.
func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

// graphqlJSON wraps data in a standard GraphQL response envelope.
func graphqlJSON(data string) string {
	return fmt.Sprintf(`{"data":%s}`, data)
}
