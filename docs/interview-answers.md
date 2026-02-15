# Interview Answers: PromptQL TUI — Read-Only CI Tests

## Project Name
`promptql-tui-readonly-tests`

## Summary
Add a comprehensive CI step to the PromptQL TUI project that tests the read-only portions of the CLI. This includes unit tests for the SDK layer (using interface-based mocks) and behavior tests for the TUI layer (Bubble Tea model/view logic). The CI pipeline runs on GitHub Actions with `go test`, `go vet`, `golangci-lint`, and `staticcheck`.

## Technical Decisions

| Decision | Choice |
|----------|--------|
| Language | Go (matches existing project, v1.24) |
| Test framework | Go standard `testing` package |
| Mocking strategy | Interface-based mocks (swap HTTP client with mock in tests) |
| CI platform | GitHub Actions |
| Linting | `golangci-lint` + `staticcheck` |
| Trigger | Push and PR to `main` |

## Feature List

### Must-Have (v1)

**SDK Unit Tests (interface-based mocks):**
1. `Projects().ListUserProjects()` — happy path + error cases
2. `Projects().GetConfig()` — happy path + error cases
3. `Projects().GetPlaygroundConfig()` — happy path + error cases
4. `Threads().List()` — happy path + error cases

**Error scenarios for each SDK method:**
- Authentication errors (invalid/expired PAT)
- Runtime errors
- Malformed API responses

**TUI Behavior Tests:**
- Config loading (from file and environment variables)
- View state transitions (setup -> projects -> threads -> chat)
- Project list rendering and navigation
- Thread list rendering and navigation
- Error state handling

**GitHub Actions CI Pipeline:**
- `go test ./...` — run all tests
- `go vet ./...` — static analysis
- `golangci-lint` — comprehensive linting
- `staticcheck` — additional static checks
- Triggered on push and pull request to `main`

### Nice-to-Have (not in scope)
- Tests for mutating operations (thread creation, message sending, API key management)
- Integration tests with real API
- Code coverage reporting/thresholds
- Tests for `Query().Ask()` / `Query().Execute()`

## Architecture Decisions

- **Interface extraction:** Introduce an interface for the HTTP transport layer in the SDK so it can be swapped with mocks during testing. This is the minimal refactoring needed to make the existing code testable.
- **Test file placement:** Tests live alongside source files following Go conventions (`*_test.go` in the same package).
- **No test fixtures directory needed:** Mock responses will be defined inline in test files for clarity.

## Constraints & Requirements

- Must not break existing functionality (no changes to runtime behavior)
- Interface extraction should be minimal and non-invasive
- Tests must run without network access (fully mocked)
- CI must work with Go 1.24
- All tests must pass before merge
