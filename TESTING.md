# Testing Guide

Quick reference for running tests in the go-hris project.

## Quick Start

```bash
# Unit tests
make test

# Integration tests (recommended for full validation)
make test-integration

# Quick smoke test (server must be running)
make test-integration-quick

# Detailed endpoint tests (server must be running)
make test-integration-endpoints
```

## Test Commands

### Unit Tests
| Command | Description |
|---------|-------------|
| `make test` | Run Go unit tests (`go test ./...`) |
| `make test-coverage` | Run tests with coverage report (generates `coverage.html`) |

### Integration Tests
| Command | Description | Server Required? |
|---------|-------------|------------------|
| `make test-integration` | Full test suite with automatic server management | ❌ No (auto-starts) |
| `make test-integration-quick` | Fast smoke test (4 core endpoints) | ✅ Yes |
| `make test-integration-endpoints` | Comprehensive endpoint tests (10 test cases) | ✅ Yes |

## Test Output Example

```
Building server...
Server built successfully
Starting server on port 8080...
Waiting for server to be ready...
Server is ready!

Running tests...
✓ Ping returns 200
✓ Login returns 200
✓ Get users returns 200
...

Total tests:  10
Passed:       10
Failed:       0

All tests passed! 🎉
```

## Development Workflow

### 1. During Development

```bash
# Terminal 1: Start server
make server

# Terminal 2: Run quick tests after changes
make test-quick
```

### 2. Before Committing

```bash
# Run full test suite
make test

# Verify build
go build ./...
```

### 3. CI/CD

The `make test` command is designed for CI/CD pipelines:
- ✅ Self-contained (builds and starts server)
- ✅ Automatic cleanup
- ✅ Returns proper exit codes
- ✅ No manual intervention needed

## More Information

For detailed documentation and CI/CD integration examples, see scripts in the `scripts/` directory.

## Test Scripts Location

```
scripts/
├── run_tests_with_server.sh  # Full automation (make test)
├── test_endpoints.sh          # Detailed tests (make test-endpoints)
├── quick_test.sh              # Smoke tests (make test-quick)
└── README.md                  # Detailed documentation
```
