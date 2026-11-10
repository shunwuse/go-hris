# Development

Run commands from the repository root. Use Go 1.26.2 or later, and have a
PostgreSQL instance available. Docker and Compose are only needed for the
container workflows. Start with the [README quick start](../README.md#-quick-start).

## Configuration

Settings are loaded from [`configs/default.env`](../configs/default.env), then
optional `configs/.env`, then `APP_` environment variables. For example,
`APP_PORT=9090 make run` overrides `port` for the API process. For a different
PostgreSQL instance, configure the `db_*` keys before applying migrations or
starting either process.

The default `use_miniredis=true` starts an embedded cache in each process.
To share a Redis instance, set `use_miniredis=false` and configure the
`redis_*` keys.

## Common commands

| Command | Purpose |
| --- | --- |
| `make` or `make help` | List Make targets. |
| `make run`, `make run-worker` | Run the API or worker locally. |
| `make fmt` | Format Go source. |
| `make modernize` | Apply x/tools modernization rewrites; review the diff afterward. |
| `make test` | Run the full Go test suite. See [Testing](#testing) for focused commands. |
| `make gen` | Regenerate Ent and Wire code and tidy modules. |
| `make docker-build` | Build and load a local Docker image. |

## Code generation

Run `make gen` after changing Ent schemas or Wire providers. The
[generation script](../scripts/gen.sh) runs Ent generation first, then Wire
for both entry points, then `go mod tidy`. Commit generated files alongside
their source changes instead of editing them by hand. If an Ent change
requires a PostgreSQL schema change, add a SQL migration as well; Ent code
generation does not migrate the database.

## Testing

Run `make test` or `go test ./...` from the repository root. The suite uses
Go's `testing` package, [testify](https://github.com/stretchr/testify), and
in-process dependencies such as `httptest` and Miniredis, so it does not
need a running database or cache. To exercise migrations and seed data
against a real database, follow the [demo data workflow](#demo-data).

| Area | Command |
| --- | --- |
| HTTP controllers only | `go test -v ./internal/http/controllers/...` |
| All HTTP tests (controllers, routing, parsing, and middleware) | `go test -v ./internal/http/...` |
| Business logic and permissions | `go test -v ./internal/services/...` |
| Infrastructure and shared utilities | `go test -v ./internal/infra/... ./internal/pkg/...` |
| Composed queries and data mapping | `go test -v ./internal/queries/... ./internal/repositories/...` |
| Scheduler and worker lifecycle | `go test -v ./internal/worker/... ./cmd/worker` |
| Seed options, naming, and counts | `go test -v ./cmd/seed` |

Generate an HTML coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Keep tests beside the code they exercise. See the
[approval controller tests](../internal/http/controllers/approval_test.go)
for HTTP examples, and update the [test doubles](../internal/mocks) when
their interfaces change.

## Migrations

```bash
go run ./cmd/migrate up
make migrate-create name=describe_the_change
```

`go run ./cmd/migrate down` rolls back **all** applied migrations, not just
the latest one. This can destroy data; use it only on a disposable database.
If `up` reports `no change`, there are no new migrations to apply; the
current command still exits nonzero in that case.

## Demo data

After applying migrations to a disposable development database:

```bash
go run ./cmd/seed
go run ./cmd/seed -users=100 -approvals=300 -seed=42
SEED_PASSWORD='local-demo-only' go run ./cmd/seed -users=2 -approvals=0
go run ./cmd/seed -users=0 -approvals=19
```

The defaults are **20 users and 50 approvals**. The command appends a batch
in one transaction and refuses to write unless `env=development`. Run one
batch at a time and don't point it at shared or production databases.

New users get unique numbered `demo.staff.*` or `demo.manager.*` usernames.
New demo accounts use `password` by default. Set `SEED_PASSWORD` to a
non-empty value of up to 72 bytes to override it for the new accounts in
that batch. Set `-approvals=0` for users only. Approval-only runs need
existing active demo staff.

Approvals use active demo staff as creators. About 10% are approved and 5%
rejected (both rounded down); the rest are pending and have no approver.
Resolved items use active demo managers or fall back to the migrated
administrator account when no demo managers exist. `-seed` controls
shuffling and assignments for the same starting database state, not database
IDs, timestamps, or password hashes.

## Docker images

See [Run with Docker](../README.md#run-with-docker) for both standalone and
Compose examples. `make docker-build` loads `go-hris:latest` for `linux/amd64`
by default and embeds the short Git revision as the version.

```bash
BUILD_MODE=debug make docker-build
PLATFORM=linux/arm64 DOCKER_IMAGE=go-hris:arm64 VERSION=1.0.0 make docker-build
```

Release builds trim paths and strip debug information. Debug builds retain
debug information and disable compiler optimizations. The image packages
both the API and worker; it does not run migrations or bundle a database.

### Connect to your own database

For the standalone Docker example, create an environment file outside the
repository with the database settings for the container. For example:

```dotenv
APP_DB_HOST=db.example.com
APP_DB_PORT=5432
APP_DB_USER=hris
APP_DB_PASSWORD=replace-me
APP_DB_NAME=hris
APP_DB_SSL_MODE=require
```

Replace the example values with your own. Use a database address reachable
from inside the container: `localhost` there refers to the container, not
the host. Apply migrations from the host to the same database before
starting the API. The host and container may need different addresses to
reach it. If you already have the prebuilt `shunwuse/go-hris:latest` image,
you can use it instead of `go-hris:latest` in the standalone command with
the same database settings.

## Troubleshooting

- **The API cannot connect to the database:** confirm the database is ready
  and check the host, port, credentials, and SSL mode. For a container,
  check that the configured address is reachable from inside it.
- **Health check returns 401 or 403:** use the configured profiler token and
  header, as shown in the [health check example](runtime.md#health-metrics-and-profiles).
- **The approval-only seed command finds no demo staff:** run a batch that
  creates demo users first.
