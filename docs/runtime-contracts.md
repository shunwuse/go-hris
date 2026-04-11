# Runtime Contracts

This document records a few implementation contracts that are easy to miss when working on the project.

## Trace ID and Idempotency

The current HTTP stack intentionally reuses `X-Trace-Id` as the idempotency key.
The goal is to avoid generating and propagating a second unique id for the same logical request.

That is a valid tradeoff, but the contract needs to stay explicit:

- Callers or proxies that retry the same logical request must preserve the same `X-Trace-Id`.
- If `X-Trace-Id` is regenerated, rewritten, or dropped, idempotency behavior changes.
- This is different from a dedicated `Idempotency-Key` header or a request fingerprinting scheme.

Related code:

- [Trace middleware](../internal/http/middlewares/trace.go)
- [Idempotency middleware](../internal/http/middlewares/idempotency.go)

## Other Runtime Notes

- `cache.Fetch` is the active cache-aside helper used by user identity reads.
- `metrics` middleware is wired into the HTTP stack and exposed at `/metrics`.
- `lock`, `lifecycle`, and `routine` are available infra components, but their runtime consumers are still limited.
