package constants

type contextKey string

const (
	// TraceID is the context key for storing trace ID.
	TraceID contextKey = "trace_id"

	// JWTClaims is the context key for storing JWT claims.
	JWTClaims contextKey = "claims"

	// DBTransaction is the context key for storing database transaction.
	DBTransaction contextKey = "db_transaction"
)
