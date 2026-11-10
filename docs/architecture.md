# 🏗️ Architecture

## 🧩 System Shape

Go HRIS is organized as a layered application:

```text
HTTP route/controller -> service -> repository -> Ent/database
                         |
                         +-> ports for cross-layer contracts
```

The API server and worker share configuration, logging, database, cache, authentication, and alerting infrastructure. They are separate binaries with separate startup and shutdown flows.

## 📂 Project Structure

The following is the current folder layout of the repository:

```text
.
├── build/
│   └── package/             # Container and packaging assets
├── cmd/                     # Application entry points
├── configs/                 # Configuration files
├── deployments/             # Docker Compose and deployment assets
├── docs/                    # Documentation and guides
├── ent/                     # Ent ORM schemas and generated code
├── internal/
│   ├── constants/            # Shared application constants
│   ├── domains/              # Domain models and claims
│   ├── dtos/                 # Request and response data structures
│   ├── errors/               # Application and domain errors
│   ├── http/                 # Router, routes, middleware, controllers, and I/O
│   ├── infra/                # Database, cache, tokens, metrics, and lifecycle tools
│   ├── mocks/                # Test doubles for interfaces
│   ├── pkg/                  # Shared internal utilities
│   ├── ports/                # Cross-layer interfaces
│   ├── queries/              # Reusable query definitions
│   ├── repositories/         # Ent-backed data access
│   ├── services/             # Business logic
│   └── worker/               # Scheduler, jobs, consumer, and handlers
├── scripts/                  # Reusable local and CI scripts
└── migrations/               # SQL migrations
```

Runtime-generated and local-only artifacts are not part of the application source layout.

### 🗺️ Planned Additions (Roadmap)

The following path is planned for future CI helper assets and is not part of the current layout:

```text
.
└── build/
    └── ci/                # Optional CI helper assets when the pipeline grows
```

## 🔌 Dependency Injection

Google Wire assembles dependencies for both binaries. Provider sets are grouped by concern, including HTTP, infrastructure, services, repositories, and worker components.

Generated injectors live in:

- `cmd/server/wire_gen.go`
- `cmd/worker/wire_gen.go`

When adding a provider, update the corresponding `wire.go` file and run `make gen`.

## 🚀 Server Composition

The server creates:

1. Application configuration and logger.
2. PostgreSQL/Ent database access.
3. Redis cache access, or embedded Miniredis when configured.
4. Authentication/token services.
5. Repositories and business services.
6. HTTP controllers, routes, and common middleware.
7. The HTTP server with graceful shutdown handling.

## ⚙️ Worker Composition

The worker starts two managed components concurrently:

- A cron scheduler. The current scheduled job cleans up expired refresh tokens daily at 02:00.
- A consumer registry. Handler registration exists, but an external message broker is not currently connected.

Both components receive cancellation through a shared context and are stopped before database and cache connections are closed.

## 🛠️ Generated Code

Ent generation must happen before Wire generation because Wire depends on the generated Ent types:

```text
go generate ./ent
-> go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/server
-> go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/worker
-> go mod tidy
```

Use `make gen` instead of editing generated files manually.
