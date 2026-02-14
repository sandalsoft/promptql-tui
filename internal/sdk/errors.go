package sdk

import "fmt"

// PromptQLError is the base error type for all PromptQL SDK errors.
type PromptQLError struct {
	Message    string
	StatusCode int
	Detail     string
}

func (e *PromptQLError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("PromptQLError (HTTP %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("PromptQLError: %s", e.Message)
}

// AuthenticationError is returned on HTTP 401 responses.
type AuthenticationError struct{ PromptQLError }

// ForbiddenError is returned on HTTP 403 responses.
type ForbiddenError struct{ PromptQLError }

// NotFoundError is returned on HTTP 404 responses.
type NotFoundError struct{ PromptQLError }

// ValidationError is returned on HTTP 422 responses.
type ValidationError struct{ PromptQLError }

// RateLimitError is returned on HTTP 429 responses.
type RateLimitError struct{ PromptQLError }

// ServerError is returned on HTTP 5xx responses.
type ServerError struct{ PromptQLError }

// newErrorForStatus returns the appropriate typed error for the given HTTP status code.
func newErrorForStatus(status int, message string) error {
	base := PromptQLError{Message: message, StatusCode: status, Detail: message}
	switch status {
	case 401:
		return &AuthenticationError{base}
	case 403:
		return &ForbiddenError{base}
	case 404:
		return &NotFoundError{base}
	case 422:
		return &ValidationError{base}
	case 429:
		return &RateLimitError{base}
	}
	if status >= 500 && status < 600 {
		return &ServerError{base}
	}
	return &base
}
