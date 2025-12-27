package response

// ========================================
// Error Response Structures
// ========================================

// ErrorResponse wraps error details
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains structured error information
type ErrorDetail struct {
	Code    string            `json:"code"`              // Domain error code
	Message string            `json:"message"`           // Human-readable message
	Details map[string]string `json:"details,omitempty"` // Optional structured context
}

// ========================================
// List Response Structures
// ========================================

// OffsetListResponse wraps data with offset-based pagination metadata
type OffsetListResponse struct {
	Data any                  `json:"data"`
	Meta OffsetPaginationMeta `json:"meta"`
}

// CursorListResponse wraps data with cursor-based pagination metadata
type CursorListResponse struct {
	Data any                  `json:"data"`
	Meta CursorPaginationMeta `json:"meta"`
}

// ========================================
// Pagination Metadata Structures
// ========================================

// OffsetPaginationMeta contains metadata for offset-based pagination
type OffsetPaginationMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

// CursorPaginationMeta contains metadata for cursor-based pagination
type CursorPaginationMeta struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}
