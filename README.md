# Human Resources Information System

## Description
This is a simple Human Resources Information System (HRIS).

## How to run
```plaintext
1. Clone this repository
2. Run `make go-migrate-up` to migrate the database
3. Run `make server` to start the server on port 8080 (default)
```

Alternatively, you can run the server using Docker in two ways
```plaintext
Build and run the server from source:
1. Clone this repository
2. Run `make docker build` to build the docker image
3. Run `make docker run` to run the docker image, the server will be available on port 8080

Run the server using a pre-built image from Docker Hub:
1. Run `docker run --rm -p 8080:8080 shunwuse/go-hris:latest`
```
**Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

[Postman Collection](https://documenter.getpostman.com/view/23207346/2sA3duEsLN)


Login with default user:
```plaintext
username: admin
password: password
```

## Testing

### Quick Start

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

### Test Commands

| Command | Description | Auto-starts Server? |
|---------|-------------|---------------------|
| `make test` | Go unit tests (`go test ./...`) | N/A |
| `make test-coverage` | Unit tests with coverage report | N/A |
| `make test-integration` | Full integration test suite | ✅ Yes |
| `make test-integration-quick` | Fast smoke test (4 core endpoints) | ❌ No |
| `make test-integration-endpoints` | Comprehensive tests (10 test cases) | ❌ No |

### Output Example

```
Building server...
Server built successfully
Starting server on port 8080...
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

### Documentation

- 📖 **[TESTING.md](TESTING.md)** - Quick testing guide

## Features
- [x] Create user
- [x] Login
- [x] Role
- [x] Permission
- [x] Approval Management
