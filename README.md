# 🚀 Go HRIS

**A Human Resources Information System (HRIS) practice project built with Go.**

*This is a side project for learning purposes and is **not intended for production use**.*

## 🛠️ Tech Stack

- **Language**: Go 1.25.7
- **Framework**: Chi (HTTP Router)
- **ORM**: Ent
- **DI**: Google Wire
- **Config**: Koanf
- **Logging**: Zap
- **Database**: PostgreSQL
- **Container**: Docker (Multi-stage build)

## 🚀 Quick Start

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/shunwuse/go-hris.git
   cd go-hris
   ```

2. **Initialize the database**
   ```bash
   make go-migrate-up
   ```

3. **Start the api server**
   ```bash
   make run
   ```
   The server will be running at `http://localhost:8080` by default.

4. **Format/Modernize the code**
   ```bash
   make fmt
   ```
   Ensures code style consistency and applies latest Go best practices.

5. **Start the worker (Optional)**
   ```bash
   make run-worker
   ```
   The worker handles background tasks like token cleanup.

### Run with Docker

#### Build from source
```bash
# Build the image (automatically injects Git Commit Hash as version)
make docker-build

# Run the container
make docker-run
```

#### Use pre-built image
```bash
docker run --rm -p 8080:8080 shunwuse/go-hris:latest
```

## 📖 API Documentation

- **Health Check**: `GET /health` (Check system status and version info)
- **Postman Collection**: [View here](https://documenter.getpostman.com/view/23207346/2sA3duEsLN)

### Default Admin Account
- **Username**: `admin`
- **Password**: `password`

## ⚙️ Configuration

The application uses a multi-stage configuration loading process:

1. **Default Settings**: Loaded from [configs/default.env](configs/default.env).
2. **Local Overrides**: Loaded from `configs/.env`.
3. **Environment Variables**: System variables with `APP_` prefix (e.g., `APP_SERVER_PORT=9090`).

To customize your local environment, we recommend creating a `.env` file in the `configs/` directory.

## 🧪 Testing Guide

We provide a comprehensive testing suite focusing on unit and layer-based testing.

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests |
| `make test-controllers` | Run only HTTP controller tests |
| `make test-services` | Run only service layer tests |
| `make test-coverage` | Generate HTML coverage report |

For detailed testing architecture and instructions, please refer to [docs/testing.md](docs/testing.md).

## 📂 Project Structure

```text
.
├── build/                # Packaging and CI related files
│   ├── ci/               # CI configurations and scripts
│   ├── package/          # Container and OS specific packages
│   └── scripts/          # Build, lint, and release scripts
├── cmd/                  # Application entry points
├── configs/              # Configuration files (.env)
├── deployments/          # System configurations and scripts
├── docs/                 # Documentation and guides
├── ent/                  # Ent ORM schemas and generated code
├── internal/             # Core business logic (Layered Architecture)
│   ├── constants/        # Application-wide constants
│   ├── domains/          # Pure domain models (Business entities)
│   ├── dtos/             # Data Transfer Objects (Request/Response)
│   ├── errors/           # Custom error definitions
│   ├── http/             # HTTP controllers, middleware and routing
│   ├── infra/            # Infrastructure
│   ├── pkg/              # Internal helper packages
│   ├── ports/            # Interfaces (Service & Repository definitions)
│   ├── repositories/     # Repository implementations (Ent)
│   └── services/         # Business logic implementations
└── migrations/           # SQL migration files (golang-migrate)
```

## 📜 Development Commands (Makefile)

Use `make` or `make help` to see all available commands. Key commands:

- `make gen`: Generate Wire DI (Required after adding dependencies)
- `make run`: Run the API server locally
- `make run-worker`: Run the background worker locally
- `make build`: Build local binary
- `make test`: Run all unit tests
- `make migrate-create name=xxx`: Create a new database migration file
- `make docker-build`: Build production Docker image

---

## ✨ Features
- [x] User Management (CRUD)
- [x] JWT Authentication & Authorization
- [x] Role-Based Access Control (RBAC)
- [x] Approval Workflow Management
- [x] Redis-based Idempotency & Distributed Locking
- [x] Structured Logging & Health Checks
