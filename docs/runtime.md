# Runtime

## HTTP request handling

The [common middleware](../internal/http/middlewares/middlewares.go) runs in
this order: CORS → trace → metrics → request logging → recovery →
idempotency → exception alerting. Requests that reach the trace middleware
get an `X-Trace-Id` (supplied by the caller or generated) and a new
`X-Span-Id` response header. Both IDs are available in the request context.

### Trace ID and idempotency

For methods other than GET, HEAD, and OPTIONS, the
[idempotency middleware](../internal/http/middlewares/idempotency.go) can
replay a stored response. It stores the status, body, and `Content-Type`
(if present) for responses with a recorded status of 200–499. The retention
period is configured with `idempotency_expire_minutes`. A replay sets
`X-Idempotency-Hit: true`.

The middleware reuses `X-Trace-Id` as the idempotency key rather than taking
a separate header or computing a request fingerprint. Callers and proxies
must preserve the same trace ID when retrying the same logical request.
The [trace middleware](../internal/http/middlewares/trace.go) must run
first so the ID is available in the request context.

### Authentication

`POST /login` and `POST /auth/refresh` are outside the access-token
middleware; refresh requires a refresh token. Routes under `/users` and
`/approvals`, plus `POST /auth/logout` and `POST /auth/logout-all`, require
`Authorization: Bearer <access-token>`. Access-token validation includes
checking expiry, token type, and revocation. The auth service loads user
roles and permissions through the cache-aside identity reader.

## Health, metrics, and profiles

[Monitor routes](../internal/http/routes/monitor.go) require the configured
`X-Profiler-Token` header. The middleware returns 403 if the server has no
profiler token configured and 401 if the supplied token is missing or wrong.

| Endpoint | What it reports |
| --- | --- |
| `GET /health` | Runtime and version metadata plus PostgreSQL and cache status. Returns HTTP 503 when health is down. |
| `/metrics` | Prometheus request counters and duration histograms, plus Go and process metrics. Unmatched requests use the `unmatched` path label. |
| `/debug/pprof/` | Go profiling endpoints. |

With the unmodified local defaults:

```bash
curl -fsS -H 'X-Profiler-Token: profiler-token' http://localhost:8080/health
```

## Storage and alerts

- PostgreSQL is accessed through sqlx and Ent. Database startup checks
  connectivity with a five-second timeout. Use the transaction abstraction
  for multi-step writes: `WithTx` commits on success and rolls back on error
  or panic. Obtain the Ent client from the callback context.
- Redis backs cache-aside reads, access-token revocation, and idempotency.
  The [user identity query](../internal/queries/user_identity.go) uses
  `cache.Fetch` as its cache-aside helper. Embedded Miniredis is process-local
  and loses state on restart.
- The infrastructure layer also provides helpers for locks, lifecycle
  coordination, and routines. They are not what drives server and worker
  shutdown. The Redis lock helper requires a TTL of at least 100 ms at
  acquisition; HTTP idempotency is not a lock.
- HTTP 5xx responses and recovered HTTP and scheduler panics can notify the
  shared alerter, which fans out to logger and terminal providers with a
  ten-second timeout.

## Process lifecycle

The [API server](../cmd/server/server.go) handles SIGINT and SIGTERM,
gives HTTP shutdown a 30-second budget, then closes cache and database
connections. The [worker](../cmd/worker/worker.go) cancels its scheduler and
consumer on shutdown, runs their stop hooks with a 30-second context, then
closes the same dependencies. Scheduler jobs run with a system identity and
a generated trace ID.
