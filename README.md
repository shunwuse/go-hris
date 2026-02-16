# 🚀 Go HRIS

**A Human Resources Information System (HRIS) practice project built with Go.**

*This is a side project for learning purposes and is **not intended for production use**.*

## 🛠️ Tech Stack

- **Language**: Go 1.25.7
- **Framework**: Chi (HTTP Router)
- **ORM**: Ent
- **DI**: Google Wire
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

We provide a comprehensive test suite including unit and integration tests.

| Command | Description | Auto-starts Server? |
|---------|-------------|---------------------|
| `make test` | Run all unit tests | No |
| `make test-coverage` | Generate test coverage report | No |
| `make test-integration` | Run full integration test suite | ✅ Yes |
| `make test-integration-quick` | Fast smoke test (requires running server) | No |

For detailed testing instructions, please refer to [docs/testing.md](docs/testing.md).

## 📂 Project Structure

- `build/`: Packaging and CI related files
  - `package/`: Container and OS specific packages
  - `ci/`: Continuous Integration configurations and scripts
  - `scripts/`: Scripts for various build, lint, and release operations
- `cmd/`: Application entry points (Server, Migration)
- `configs/`: Configuration file templates or default configs
- `deployments/`: System configurations and scripts
- `ent/`: Generated Ent ORM schemas and code
- `internal/`: Core business logic
  - `domains/`: Domain models
  - `services/`: Business logic layer
  - `repositories/`: Data access layer
  - `http/`: Controllers and routes
- `migrations/`: SQL migration files
- `scripts/`: Helper scripts

## 📜 Development Commands (Makefile)

Type `make` or `make help` to see all available commands:

- `make gen`: Generate Wire dependency injection code
- `make build`: Build local binary
- `make clean`: Clean up build artifacts and temporary database
- `make migrate-create name=xxx`: Create a new database migration file

## ✨ Features
- [x] User Management (CRUD)
- [x] JWT Authentication & Authorization
- [x] Role-Based Access Control (RBAC)
- [x] Approval Workflow Management
- [x] Structured Logging & Health Checks
