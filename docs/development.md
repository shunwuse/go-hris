# 🛠️ Development Guide

## 📋 Prerequisites

- Go 1.26.2
- Docker and Docker Compose when using containerized dependencies
- PostgreSQL for a local external database, unless an existing instance is available

Local Redis can be replaced by embedded Miniredis through `configs/default.env`.

## ⚙️ Configuration

Create a local configuration file when overrides are needed:

```bash
cp configs/default.env configs/.env
```

Configuration is loaded from defaults, then `configs/.env`, then `APP_` environment variables. For example:

```bash
APP_PORT=9090
```

Use the keys defined in `configs/default.env` for the available settings.

## 📜 Common Commands

| Command | Purpose |
|---------|---------|
| `make run` | Run the API server |
| `make run-worker` | Run the worker |
| `make test` | Run the full Go test suite |
| `make fmt` | Format Go source files |
| `make modernize` | Apply Go modernization rewrites |
| `make gen` | Generate Ent and Wire code, then tidy modules |
| `make migrate-create name=example` | Create a migration file |
| `make docker-build` | Build the Docker image |

Use `make help` to list the available targets.

## 🧬 Code Generation

Run `make gen` after changing Ent schemas or Wire providers.

Generated Ent code and Wire injectors are part of the repository. Do not hand-edit generated files; update their source schema or provider definitions and regenerate them.

## 🗄️ Database

Start the PostgreSQL service from Docker Compose when needed:

```bash
docker compose -f deployments/docker-compose.yaml up -d postgres
```

Apply migrations with:

```bash
go run ./cmd/migrate up
```

Create a new migration with:

```bash
make migrate-create name=describe_the_change
```

## 🌱 Demo Data

The seed command creates demo users and approvals for pagination and workflow demonstrations:

```bash
go run ./cmd/seed
```

By default, the command adds 50 demo users and 100 demo approvals.

The record counts can be customized:

```bash
go run ./cmd/seed --users 100 --approvals 300
```

Demo users are created with manager and staff roles, and use `password` as
their password. Approval creators and statuses are shuffled, while most
approvals remain pending. Run the command again to add another batch of records.

## 🐳 Docker Build

Build the image from the repository root:

```bash
make docker-build
```

The default image is `go-hris:latest`, built for `linux/amd64` and loaded into the local Docker engine. The wrapper accepts these environment overrides:

```bash
BUILD_MODE=debug make docker-build
PLATFORM=linux/arm64 DOCKER_IMAGE=go-hris:arm64 make docker-build
VERSION=1.0.0 make docker-build
```

`release` builds use trimmed paths and linker size flags. `debug` builds retain debug information and disable compiler optimizations for easier debugging.

Run the resulting image with:

```bash
docker run --rm -p 8080:8080 go-hris:latest
```

## 🧪 Focused Tests

Run tests for one layer when iterating:

```bash
go test -v ./internal/http/controllers/...
go test -v ./internal/services/...
go test -v ./internal/infra/...
go test -v ./internal/http/...
go test -v ./internal/queries/... ./internal/repositories/...
go test -v ./internal/worker/... ./cmd/worker/...
```

For test architecture, mocks, and coverage, see [Testing Guide](testing.md).
