package domains

// OffsetQuery represents the parameters for offset-based pagination.
type OffsetQuery struct {
	Page    int `schema:"page"`
	PerPage int `schema:"per_page"`
}

// Normalize sets default values for OffsetQuery.
func (q *OffsetQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 10
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
}

// Offset calculates the database offset.
func (q OffsetQuery) Offset() int {
	return (q.Page - 1) * q.PerPage
}

// OffsetResult represents the result of an offset-based pagination.
type OffsetResult[T any] struct {
	Items      []T
	TotalCount int
	TotalPage  int
}

// CursorQuery represents the parameters for cursor-based pagination.
type CursorQuery struct {
	Cursor string `schema:"cursor"`
	Limit  int    `schema:"limit"`
}

// Normalize sets default values for CursorQuery.
func (q *CursorQuery) Normalize() {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
}

// CursorResult represents the result of a cursor-based pagination.
type CursorResult[T any] struct {
	Items      []T
	NextCursor string
	HasMore    bool
}

type ApprovalCursor struct {
	ID uint `json:"id"`
}
