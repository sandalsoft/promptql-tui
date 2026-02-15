package sdk

import (
	"errors"
	"net/http"
	"testing"
)

func TestListThreads_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		body := graphqlJSON(`{
			"getThreads": [
				{
					"thread_id": "t-1",
					"title": "First thread",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-02T00:00:00Z",
					"project_id": "proj-1",
					"build_id": "build-1",
					"user_id": "user-1",
					"visibility": "private"
				},
				{
					"thread_id": "t-2",
					"title": "Second thread",
					"created_at": "2025-01-03T00:00:00Z",
					"updated_at": "2025-01-04T00:00:00Z",
					"project_id": "proj-1",
					"build_id": "build-2",
					"user_id": "user-1",
					"visibility": "public"
				}
			]
		}`)
		return jsonResponse(200, body), nil
	})

	threads, err := client.Threads().List("proj-1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(threads) != 2 {
		t.Fatalf("expected 2 threads, got %d", len(threads))
	}

	if threads[0].ThreadID != "t-1" {
		t.Errorf("expected ThreadID 't-1', got %q", threads[0].ThreadID)
	}
	if threads[0].Title != "First thread" {
		t.Errorf("expected Title 'First thread', got %q", threads[0].Title)
	}
	if threads[0].ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got %q", threads[0].ProjectID)
	}
	if threads[0].Visibility != "private" {
		t.Errorf("expected Visibility 'private', got %q", threads[0].Visibility)
	}

	if threads[1].ThreadID != "t-2" {
		t.Errorf("expected ThreadID 't-2', got %q", threads[1].ThreadID)
	}
	if threads[1].Visibility != "public" {
		t.Errorf("expected Visibility 'public', got %q", threads[1].Visibility)
	}
}

func TestListThreads_Empty(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		body := graphqlJSON(`{"getThreads": []}`)
		return jsonResponse(200, body), nil
	})

	threads, err := client.Threads().List("proj-1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if threads == nil {
		t.Fatal("expected non-nil slice, got nil")
	}
	if len(threads) != 0 {
		t.Errorf("expected 0 threads, got %d", len(threads))
	}
}

func TestListThreads_AuthError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(401, `{"message":"invalid token"}`), nil
	})

	_, err := client.Threads().List("proj-1", "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var authErr *AuthenticationError
	if !errors.As(err, &authErr) {
		t.Errorf("expected *AuthenticationError, got %T: %v", err, err)
	}
}

func TestListThreads_RuntimeError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection reset by peer")
	})

	_, err := client.Threads().List("proj-1", "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !containsString(err.Error(), "connection reset by peer") {
		t.Errorf("expected error to contain 'connection reset by peer', got: %v", err)
	}
}

func TestListThreads_MalformedResponse(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{totally broken`), nil
	})

	_, err := client.Threads().List("proj-1", "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
