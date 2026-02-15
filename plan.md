# Implementation Plan: PromptQL TUI Read-Only CI Tests

## Summary
Add comprehensive unit tests for the read-only portions of the PromptQL TUI CLI, covering both the SDK layer (4 methods with interface-based HTTP mocks) and the TUI layer (config loading, view state transitions, list rendering, error handling). Set up a GitHub Actions CI pipeline that runs `go test`, `go vet`, `golangci-lint`, and `staticcheck` on every push and pull request.

## Technical Stack
- **Language**: Go 1.24
- **Test framework**: Go standard `testing` package
- **Mocking**: Custom `http.RoundTripper` mock (injected via `ClientOptions.HTTPClient`)
- **TUI testing**: Direct Bubble Tea `Model.Update()` calls with constructed messages
- **CI**: GitHub Actions with `golangci-lint` and `staticcheck`
- **Linting**: `golangci-lint` + `staticcheck`

## Phases

### Phase 1: Test Infrastructure — Mock HTTP Transport
- [x] Step 1.1: Create `internal/sdk/mock_transport_test.go` — a reusable `mockRoundTripper` that implements `http.RoundTripper`. It accepts a function `func(*http.Request) (*http.Response, error)` so each test can configure the exact HTTP response (status code, body, headers). Include a `newTestClient(fn)` helper that creates an `sdk.Client` with the mock transport injected via `ClientOptions.HTTPClient`. Done when: file compiles with `go build ./internal/sdk/...`.

### Phase 2: SDK Unit Tests — Projects
- [x] Step 2.1: Create `internal/sdk/projects_test.go` with tests for `ListUserProjects()`:
  - `TestListUserProjects_Success` — mock returns valid GraphQL response with 2 projects (one with build FQDN, one without). Assert correct count, field mapping (Name, DDNProjectID, BuildFQDN).
  - `TestListUserProjects_AuthError` — mock returns HTTP 401 with JSON body `{"message":"invalid token"}`. Assert error is `*AuthenticationError`.
  - `TestListUserProjects_RuntimeError` — mock transport returns a Go error (e.g., connection refused). Assert error wraps the transport error.
  - `TestListUserProjects_MalformedResponse` — mock returns HTTP 200 with invalid JSON body. Assert error is returned (not panic).
  Done when: `go test ./internal/sdk/ -run TestListUserProjects -v` passes all 4 tests.

- [x] Step 2.2: Add tests for `GetConfig()` to `internal/sdk/projects_test.go`:
  - `TestGetConfig_Success` — mock returns GraphQL response with `promptQlEnabled: true, playgroundEnabled: false`. Assert fields match.
  - `TestGetConfig_AuthError` — HTTP 401 → `*AuthenticationError`.
  - `TestGetConfig_RuntimeError` — transport error.
  - `TestGetConfig_MalformedResponse` — invalid JSON → error.
  Done when: `go test ./internal/sdk/ -run TestGetConfig -v` passes all 4 tests.

- [x] Step 2.3: Add tests for `GetPlaygroundConfig()` to `internal/sdk/projects_test.go`:
  - `TestGetPlaygroundConfig_Success` — mock returns full PlaygroundConfig with all fields populated (including pointer fields like `*bool`, `*int`, and `map[string]interface{}` featureFlags). Assert all fields.
  - `TestGetPlaygroundConfig_AuthError` — HTTP 401 → `*AuthenticationError`.
  - `TestGetPlaygroundConfig_RuntimeError` — transport error.
  - `TestGetPlaygroundConfig_MalformedResponse` — invalid JSON → error.
  Done when: `go test ./internal/sdk/ -run TestGetPlaygroundConfig -v` passes all 4 tests.

### Phase 3: SDK Unit Tests — Threads
- [x] Step 3.1: Create `internal/sdk/threads_test.go` with tests for `Threads().List()`:
  - `TestListThreads_Success` — mock returns GraphQL response with 2 threads. Assert correct count, ThreadID, Title, ProjectID, etc.
  - `TestListThreads_Empty` — mock returns empty array. Assert zero-length slice (not nil).
  - `TestListThreads_AuthError` — HTTP 401 → `*AuthenticationError`.
  - `TestListThreads_RuntimeError` — transport error.
  - `TestListThreads_MalformedResponse` — invalid JSON → error.
  Done when: `go test ./internal/sdk/ -run TestListThreads -v` passes all 5 tests.

### Phase 4: Config Tests
- [x] Step 4.1: Create `internal/config/config_test.go`:
  - `TestLoad_NoFile` — Load from a non-existent path returns zero-value Config (no error).
  - `TestLoad_ValidFile` — Write a temp config.json, load it, assert all fields.
  - `TestLoad_MalformedJSON` — Write invalid JSON to temp file, assert error.
  - `TestHasCredentials_WithPAT` — Config with PAT returns true.
  - `TestHasCredentials_WithoutPAT` — Empty Config returns false.
  - `TestSave_RoundTrip` — Save a Config, Load it back, assert equality.
  Done when: `go test ./internal/config/ -v` passes all 6 tests.

### Phase 5: TUI Behavior Tests
- [x] Step 5.1: Create `internal/tui/app_test.go` with view state transition tests:
  - `TestNew_WithCredentials` — Config with PAT starts at `viewProjects` with loading=true.
  - `TestNew_WithoutCredentials` — Empty Config starts at `viewSetup`.
  - `TestEsc_FromProjects` — Sending esc key in projects view transitions to setup view.
  - `TestEsc_FromThreads` — Sending esc key in threads view transitions to projects view.
  - `TestEsc_FromChat` — Sending esc key in chat view transitions to threads view.
  Done when: `go test ./internal/tui/ -run "TestNew|TestEsc" -v` passes all 5 tests.

- [x] Step 5.2: Add project list behavior tests to `internal/tui/app_test.go`:
  - `TestProjectsLoaded` — Sending `projectsLoadedMsg` with projects populates list and clears loading.
  - `TestProjectNavigation` — Sending j/k keys moves `projectCursor` up/down with bounds checking.
  - `TestProjectsView_Error` — Sending `errMsg` sets error state, `View()` output contains error text.
  - `TestProjectsView_Empty` — Empty project list renders "No projects found" text.
  Done when: `go test ./internal/tui/ -run TestProject -v` passes all 4 tests.

- [x] Step 5.3: Add thread list behavior tests to `internal/tui/app_test.go`:
  - `TestThreadsLoaded` — Sending `threadsLoadedMsg` populates threads list and clears loading.
  - `TestThreadNavigation` — j/k keys move `threadCursor` (accounting for "New Thread" at index 0).
  - `TestThreadsView_Error` — `errMsg` displays error in threads view.
  - `TestNewThreadSelection` — Enter on cursor=0 transitions to chat view with empty threadID.
  Done when: `go test ./internal/tui/ -run TestThread -v` passes all 4 tests.

### Phase 6: GitHub Actions CI Pipeline
- [x] Step 6.1: Create `.github/workflows/test.yml`:
  - Trigger on push and pull_request to `main`.
  - Job `test` runs on `ubuntu-latest` with Go 1.24.
  - Steps: checkout, setup-go with caching, `go vet ./...`, `go test -v -race ./...`.
  - Job `lint` runs on `ubuntu-latest`.
  - Steps: checkout, setup-go, `golangci-lint` (via `golangci/golangci-lint-action`), `staticcheck` (install + run).
  Done when: file exists and is valid YAML.

### Phase 7: Validation
- [x] Step 7.1: Run `go vet ./...` and `go test -v -race ./...` locally. All tests must pass with no vet warnings. Fix any issues found.
- [x] Step 7.2: Run `go build ./...` to verify no compilation errors in production code. Verify test count matches expectations (~30+ tests total).
