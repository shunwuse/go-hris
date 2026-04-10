# Testing Guide

This guide describes how to run tests for the go-hris project. We have transitioned from legacy shell-based integration tests to a comprehensive Go-native unit testing suite.

## Test Architecture

We use Go's native `testing` package along with:
- **[testify](https://github.com/stretchr/testify)**: For assertions and mocking.
- **net/http/httptest**: For testing HTTP controllers without starting a real server.

The test suite is divided into three main layers:

### 1. Infrastructure & Pkg Tests
Located in `internal/infra/*/*_test.go` and `internal/pkg/*/*_test.go`.
These tests verify:
- Core utilities (Config, Logger, Pagination).
- Infrastructure components (Idempotency, Cache, Redis Locking).
- We use **Miniredis** for testing Redis-dependent code without an external server.

### 2. Controller Tests
Located in `internal/http/controllers/*_test.go`.
These tests verify:
- HTTP request routing and validation.
- Response Status Codes and JSON structure.
- Correct interaction with the Service layer via mocks.

### 2. Service Tests
Located in `internal/services/*_test.go`.
These tests verify:
- Business logic in isolation.
- Integration with the Repository layer (via mocks).
- Error handling and complex edge cases.

## Running Tests

Use `go test` directly for focused runs. The Makefile keeps only the full-suite helper.

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests in the project. |
| `go test -v ./internal/http/controllers/...` | Run only the HTTP controller tests. |
| `go test -v ./internal/services/...` | Run only the service layer logic tests. |
| `go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html` | Run all tests and generate a `coverage.html` report. |

### Example: Running specific tests
```bash
# Run all tests
make test

# Run only controller tests
go test -v ./internal/http/controllers/...

# Run only service tests
go test -v ./internal/services/...
```

## Writing New Tests

### Mocks
Mocks are generated or manually defined in `internal/mocks`. When adding a new service or repository interface, ensure you update the corresponding mock.

Example of a table-driven test with mocking:
```go
func TestController(t *testing.T) {
    tests := []struct {
        name           string
        mockSetup      func(*mocks.MockService)
        expectedStatus int
    }{
        {
            name: "Success case",
            mockSetup: func(m *mocks.MockService) {
                m.On("DoWork").Return(nil)
            },
            expectedStatus: http.StatusOK,
        },
    }
    // ...
}
```

## Legacy Integration Tests
The project previously used shell-based integration tests (`scripts/*.sh`). These have been phased out in favor of Go-native unit tests using `httptest` and component-level testing with **Miniredis** to ensure a faster and more reliable CI/CD pipeline.
