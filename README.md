# 🚀 Go HRIS

**A Human Resources Information System (HRIS) practice project built with Go.**

*This is a side project for learning purposes and is **not intended for production use**.*

## 🛠️ Tech Stack

- **Language**: Go 1.26.2
- **HTTP Router**: Chi
- **ORM**: Ent
- **Dependency Injection**: Google Wire
- **Configuration**: Koanf
- **Logging**: Zap
- **Database**: PostgreSQL
- **Cache**: Redis
- **Container**: Docker (multi-stage build)

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/shunwuse/go-hris.git
   cd go-hris
   ```

2. **Configure local overrides (Optional)** <br>
   Copy the defaults and edit any settings you need to override:
   ```bash
   cp configs/default.env configs/.env
   ```

3. **Prepare the database (Optional)** <br>
   Use your own database or start one with Compose:
   ```bash
   docker compose -f deployments/docker-compose.yaml up -d postgres
   ```

4. **Apply migrations** <br>
   Run from the repository root once the database is ready:
   ```bash
   go run ./cmd/migrate up
   ```

5. **Start the API server**
   ```bash
   make run
   ```
   The API listens at `http://localhost:8080` by default.

6. **Start the worker (Optional)** <br>
   Run in a separate terminal:
   ```bash
   make run-worker
   ```
   The worker handles background tasks like token cleanup.

### Run with Docker

You can run the API in Docker with either your own database or the database provided by Compose. Apply migrations to the database you choose before starting the API.

**With your own database:** [Set the container's database connection](docs/development.md#connect-to-your-own-database), then run just the API:

```bash
make docker-build
docker run --rm -p 8080:8080 --env-file /path/to/api.env go-hris:latest
```

**With the Compose database:** Follow the database and migration steps above with the default settings, then start the API on the same network:

```bash
docker compose -f deployments/docker-compose.yaml up -d --build api
```

## 📖 Try the API

- [API guide](docs/api.md): sign in with the local demo account (`admin` / `password`) and try the list endpoints.
- [Health check](docs/runtime.md#health-metrics-and-profiles): check the API and its dependencies.
- [Demo data](docs/development.md#demo-data): populate the user and approval lists.

## ✨ Features

- Create, view, and update users
- Token-based authentication and role-based authorization
- Approval workflow and cursor pagination
- Request idempotency and a distributed lock helper
- Structured logging, health checks, and metrics

## 📚 Project Documentation

| Guide | Covers |
| --- | --- |
| [Development](docs/development.md) | Configuration, migrations, demo data, testing, and Docker |
| [Architecture](docs/architecture.md) | Directory layout, package responsibilities, and process composition |
| [Runtime](docs/runtime.md) | HTTP behavior, monitoring, storage, and shutdown |
