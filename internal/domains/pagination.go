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
