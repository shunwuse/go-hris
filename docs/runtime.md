# ⚙️ Runtime Guide

This document describes behavior that affects API callers, operators, and contributors during local or deployed execution.

## 🔗 HTTP Middleware

Common middleware is applied in this order:

1. CORS policy.
2. Trace ID and span ID propagation.
3. Prometheus request metrics.
4. Structured request logging.
5. Panic recovery.
6. Idempotency response handling.
7. 5xx exception alerting.

The order matters. Idempotency reuses the trace ID from the request context, so trace propagation must happen first. See [runtime contracts](runtime-contracts.md) for the full constraint.

### Trace Headers

The server accepts `X-Trace-Id` when supplied, otherwise it generates one. Every response includes:

- `X-Trace-Id`
- `X-Span-Id`

The IDs are also attached to the request context and included in contextual logs and alerts.

### Idempotency

State-changing requests are eligible for idempotency handling. `GET`, `HEAD`, and `OPTIONS` requests bypass it. A repeated request with the same trace ID can receive the cached status, headers, and body while the cache entry is valid.

Only responses with status codes from 200 through 499 are cached. 5xx responses are intentionally left uncached so transient failures can be retried.

### Authentication

Protected routes require an `Authorization: Bearer <token>` header. The JWT middleware validates the access token, checks its token type and expiration, and checks Redis for a blacklisted token ID.

The current route groups are:

- Public `POST /login`.
- Refresh `POST /auth/refresh`.
- Protected auth actions under `POST /auth/logout` and `POST /auth/logout-all`.
- Protected user routes under `/users`.
- Protected approval routes under `/approvals`.

## 📊 Monitoring and Diagnostics

The monitor routes are protected by the configured profiler token:

- `GET /health`: application, database, and cache health information.
- `GET /metrics`: Prometheus metrics, including request counters and duration histograms.
- `GET /debug/pprof/*`: Go runtime profiling endpoints.

Send the configured token using `X-Profiler-Token`. If no profiler token is configured, these routes return `403 Forbidden`.

The metrics registry also exposes standard Go and process collectors. Unmatched routes are recorded with the `unmatched` path label.

## 🛠️ Infrastructure Behavior

### Database

The application uses PostgreSQL through `sqlx` and Ent. Database creation performs a connectivity check with a five-second timeout. Repository code should use the database abstraction and the transaction port rather than opening independent connections.

`WithTx` commits when the callback succeeds, rolls back when it returns an error, and rolls back before re-panicking when the callback panics.

### Cache and Locks

Redis provides cache storage, token blacklist storage, idempotency records, and distributed locks. Local development can use embedded Miniredis by setting `use_miniredis=true` in the configuration.

The distributed lock implementation enforces a minimum lock TTL of 100 milliseconds and supports refresh and release operations.

### Alerts

HTTP 5xx responses and recovered scheduler panics are sent through the shared alerter. The current alerter fan-outs to logger and terminal providers in parallel with a ten-second timeout.

## 🔄 Process Lifecycle

The server handles `SIGINT` and `SIGTERM`, stops accepting HTTP requests gracefully, and then closes cache and database connections. The shutdown timeout is 30 seconds.

The worker handles the same signals, cancels the scheduler and consumer context, waits for them to stop, and then closes cache and database connections. Scheduler jobs run with a system identity and a generated trace ID.

## ⚙️ Configuration

Configuration is loaded in layers:

1. Defaults from `configs/default.env`.
2. Optional local overrides from `configs/.env`.
3. Environment variables using the `APP_` prefix.

Important runtime settings include service port, profiler token, idempotency expiration, PostgreSQL connection values, JWT expiration values, Redis connection values, and `use_miniredis`.
