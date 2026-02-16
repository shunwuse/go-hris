# Go HRIS API Server Engineering Roadmap (The "Square")

This roadmap focuses on strengthening the API server (the "Square")—ensuring the internal logic, performance, and reliability of the Go program itself are production-grade.

---

## Phase 1: Observability & Diagnostics (單體觀測與除錯)

### 1. Contextual Logging Decorator [In-Progress]
*   **File:** `internal/pkg/logger/context.go`
*   **Description:** A helper (e.g., `logger.L().WithContext(ctx)`) that automatically extracts `trace_id`, `span_id`, and `user_id` from the context and attaches them as structured fields to every log entry.
*   **Detailed Goal:** Adopt OpenTelemetry naming standards. Currently, `trace_id` and `span_id` are implemented; `user_id` extraction from Auth context is pending. Ensure professional correlation between total journey, local execution, and actor.

### 2. Request Metadata Enhancer [Won't Do]
*   **File:** `internal/infra/metadata/metadata.go`
*   **Description:** Middleware that parses the incoming `http.Request` for client metadata (Real IP via `X-Forwarded-For`, User-Agent, Device Type, and Geo-location based on IP).
*   **Detailed Goal:** Provide rich context for security auditing, analytics, and debugging without manual parsing in every handler.

### 3. Global Exception Handler & Messenger [Completed]
*   **File:** `internal/http/middlewares/exception.go`
*   **Description:** An extension to the recovery middleware (or a coordinated middleware) that, upon encountering a critical 5xx error or panic, not only logs it but also pushes a formatted alert via `IAlerter`.
*   **Detailed Goal:** Reduce Mean Time to Detection (MTTD) of production crashes by including the error stack trace and request context in the alert. Implemented via dual-middleware strategy (Recovery & Exception).

### 4. Request Profiler On-Demand [Completed]
*   **File:** `internal/http/middlewares/profiler.go`
*   **Description:** A mechanism to enable Go's `pprof` (CPU, Memory, Block) for a specific request triggered by a secure admin header.
*   **Detailed Goal:** Diagnose performance bottlenecks (e.g., high CPU usage in a specific calculation) in live production environments without restarting the server. Secured via `X-Profiler-Token`.

---

## Phase 2: Single-Node Performance & Resilience (單機性能與彈性)

### 5. SingleFlight (Cache Protection) [Completed]
*   **File:** `internal/infra/cache/cache.go`
*   **Description:** Integration of `golang.org/x/sync/singleflight` into the `Fetch` helper.
*   **Detailed Goal:** Prevent "Cache Stampede" (Thundering Herd) where multiple concurrent requests for the same expired key all hit the database. Only one request will fetch from DB while others wait and share the result.

### 6. Response Compression (Gzip/Brotli) [Won't Do]
*   **File:** `internal/http/middlewares/middlewares.go`
*   **Description:** Middleware that automatically compresses the response body based on the client's `Accept-Encoding` header.
*   **Detailed Goal:** Significantly reduce egress bandwidth usage and improve page load times for clients, especially for large JSON responses.

### 7. Resilience Wrapper (Retry/Timeout/Circuit Breaker) [Not Started]
*   **File:** `internal/infra/resilience/resilience.go`
*   **Description:** A standardized executor (likely using a library like `resilience4go` or custom logic) to wrap external IO calls.
*   **Detailed Goal:** Shield the system from cascading failures. If an external service (e.g., Email API) is down, the Circuit Breaker prevents the server from hanging and provides a fallback.

### 8. Query Cost Analyzer / Guardian [Won't Do]
*   **File:** `internal/infra/query/query_analyzer.go`
*   **Description:** A middleware or Ent hook that analyzes SQL queries before execution (using `EXPLAIN`) to estimate their "cost" or execution time.
*   **Detailed Goal:** Systematically block "Query of Death" scenarios that would otherwise lock database tables or exhaust CPU in production.

---

## Phase 3: Resource Lifecycle & Config (資源生命週期與配置)

### 9. Lifecycle Manager (Graceful Shutdown Registry) [Completed]
*   **File:** `internal/infra/lifecycle/lifecycle.go`
*   **Description:** A central registry created via `lifecycle.New()` that allows any component to register its `Shutdown()` function.
*   **Detailed Goal:** Ensure all components (DB, Redis, MQ) are closed in the correct order (e.g., stop accepting requests first, then close DB) to avoid data corruption during restarts.

### 10. Graceful Job Runner (RoutineGroup) [Completed]
*   **File:** `internal/infra/routine/routine_group.go`
*   **Description:** A background task coordinator created via `routine.New()` using `sync.WaitGroup` to track active long-running goroutines.
*   **Detailed Goal:** Ensure the server waits for background tasks (like a daily report generation) to finish or reach a checkpoint before the process exits.

### 11. Dynamic Config (Hot-Reload) [Completed]
*   **File:** `internal/infra/config/config.go`
*   **Description:** Enhancement of the `config.Get()` to support reloading via `config.Reload()` from disk (or environment) without restarting the application, using atomic storage for thread-safety.
*   **Detailed Goal:** Change log levels or toggle feature flags instantly without downtime.

---

## Phase 4: Distributed System Reliability (分散式系統可靠性)

### 12. Schema Migrator Wrapper [Not Started]
*   **File:** `internal/infra/migrate/migrator.go`
*   **Description:** A wrapper around `golang-migrate` that performs sanity checks (e.g., checksum validation) at startup.
*   **Detailed Goal:** Prevent the application from starting if the database schema is inconsistent with the code version, preventing runtime "column not found" errors.

### 13. Health Registry [Completed: Service-Driven Approach]
*   **File:** `internal/services/monitor.go`
*   **Description:** A system where infrastructure status (DB, Redis) is checked. Current implementation uses `MonitorService` to orchestrate calls to `MonitorRepository`.
*   **Detailed Goal:** Provide accurate readiness/liveness data to Kubernetes or monitoring tools.
*   **Implementation Note:** Current approach is Service-Driven. For a more decoupled "Registry Pattern", consider refactoring each component to self-register its checker to a central registry in `internal/infra/health`.
*   **Status Detail**: Fully functional with deep checks for DB and Redis.

### 14. Idempotency Manager [Completed]
*   **File:** `internal/infra/idempotency/idempotency.go`
*   **Description:** A Redis-backed system that tracks a unique `X-Trace-Id` (reusing Trace ID context) provided by the client for state-changing operations.
*   **Detailed Goal:** Guarantee that the same request processed twice only results in one side effect. Uses Chi's native `ResponseWriter` with `Tee` for efficient response capture and supports per-endpoint custom TTL.

### 15. Distributed Lock (Redis-based) [Completed]
*   **File:** `internal/infra/lock/locker.go`
*   **Description:** A wrapper (using `bsm/redislock`) to provide mutual exclusion across multiple server instances.
*   **Detailed Goal:** Ensure that scheduled tasks or critical unique operations (like batch processing) are only executed by one instance at a time in a cluster.

### 16. Unified Cache Synchronizer [Not Started]
*   **File:** `internal/infra/cache/cache_sync.go`
*   **Description:** An Ent hook or Transactor extension that triggers cache invalidation across the Redis cluster upon successful database updates.
*   **Detailed Goal:** Solve the "Stale Cache" problem by ensuring that as soon as a DB record changes, its cached version is cleared or updated globally.

### 17. Task Scheduler (Distributed Crontab)
*   **File:** `internal/infra/scheduler/scheduler.go`
*   **Description:** A robust scheduler (like `gocron` or `asynq`) that manages recurring or delayed jobs across the server fleet.
*   **Detailed Goal:** Handle long-term scheduling (e.g., "send anniversary notification in 3 months") and ensure exactly-once execution.

---

## Phase 5: Engineering Excellence & Safety (工程卓越與安全)

### 18. Internal Event Dispatcher
*   **File:** `internal/infra/event/event_bus.go`
*   **Description:** A publisher-subscriber system (In-memory) for domain events (e.g., `UserCreated`).
*   **Detailed Goal:** Decouple primary business logic from side effects. The `UserService` only creates the user; listeners handle sending emails or updating metrics asynchronously.

### 19. Feature Toggler (Feature Flags)
*   **File:** `internal/infra/feature/feature_flag.go`
*   **Description:** A system to manage feature activation based on conditions (Percentage of users, specific UserIDs, or Global toggles).
*   **Detailed Goal:** Decouple Deployment from Release. Lower the risk of new features by enabling them gradually in production.

### 20. Universal Storage Wrapper
*   **File:** `internal/infra/storage/storage.go`
*   **Description:** An abstraction layer for file storage providing a unified interface for `LocalFileSystem`, `AWS S3`, or `Google Cloud Storage`.
*   **Detailed Goal:** Allow the developer to code against a single interface and switch storage backends via configuration.

### 21. PII Masker / Data Anonymizer
*   **File:** `internal/infra/security/pii_masker.go`
*   **Description:** A processing layer (Logger hook or Middleware) that detects patterns like Emails, Phone numbers, or IDs and masks them (e.g., `te***@example.com`).
*   **Detailed Goal:** Ensure compliance with data privacy laws (GDPR/SOC2) by preventing sensitive data from being stored in logs or accidental leaks in APIs.

### 22. Fine-Grained Access Control (Casbin/PBAC)
*   **File:** `internal/infra/authz/casbin.go`
*   **Description:** Integration of Casbin as the central Authorization engine. Supports PERM (Policy, Effect, Request, Matcher) models for RBAC and ABAC.
*   **Detailed Goal:** Decouple authorization logic from business services. Enable complex rules (e.g., "Managers can only update subordinates' records during working hours") without modifying application code. Integrated with Ent via `ent-adapter`.
