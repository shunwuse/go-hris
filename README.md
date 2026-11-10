# 🚀 Go HRIS

**A Human Resources Information System (HRIS) practice project built with Go.**

*This is a side project for learning purposes and is **not intended for production use**.*

## 🛠️ Tech Stack

- **Language**: Go 1.26.2
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

2. **Create local configuration (Optional)**
   ```bash
   cp configs/default.env configs/.env
   ```

3. **Prepare a database (Optional)** <br>
   Use an existing database, or start the database defined in Docker Compose:
   ```bash
   docker compose -f deployments/docker-compose.yaml up -d postgres
   ```

4. **Initialize the database**
   ```bash
   go run ./cmd/migrate up
   ```

5. **Start the api server**
   ```bash
   make run
   ```
   The server will be running at `http://localhost:8080` by default.

6. **Start the worker (Optional)**
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
docker run --rm -p 8080:8080 go-hris:latest
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
3. **Environment Variables**: System variables with `APP_` prefix (e.g., `APP_PORT=9090`).

To customize your local environment, we recommend creating a `.env` file in the `configs/` directory.

## 📂 Project Structure

The project structure, package responsibilities, and server/worker composition are documented in [docs/architecture.md](docs/architecture.md).

## 📜 Development Commands (Makefile)

Use `make` or `make help` to see the common local commands. The Makefile stays intentionally small so the most common workflows are easy to remember.

Key commands:

- `make gen`: Generate Wire DI and Ent code (required after adding dependencies or schemas)
- `make run`: Run the API server locally
- `make run-worker`: Run the background worker locally
- `make fmt`: Format the code with Go's standard formatter
- `make modernize`: Apply modern Go modernization rewrites for newer Go idioms
- `make test`: Run all unit tests
- `make migrate-create name=xxx`: Create a new database migration file
- `make docker-build`: Build the production Docker image

For detailed development setup, database, seed data, Docker, and focused test instructions, please refer to [docs/development.md](docs/development.md).

### 🌱 Demo Data

Generate demo users and approvals for pagination demonstrations:

```bash
go run ./cmd/seed
```

For record count options and more details, please refer to [docs/development.md](docs/development.md).

## ✨ Features
- [x] User Management (CRUD)
- [x] JWT Authentication & Authorization
- [x] Role-Based Access Control (RBAC)
- [x] Approval Workflow Management
- [x] Redis-based Idempotency & Distributed Locking
- [x] Structured Logging & Health Checks

## 📚 Runtime Contracts

Some implementation details are intentionally shared across the project and documented for contributors in [docs/runtime-contracts.md](docs/runtime-contracts.md).

For HTTP middleware, infrastructure components, monitoring, authentication, and process lifecycle behavior, please refer to [docs/runtime.md](docs/runtime.md).
