package query

type PaginationParams struct {
	Page  int
	Limit int
}

func NewPaginationParams(page, limit *int) PaginationParams {
	// defaults
	p := 1
	l := 100

	if page != nil && *page > 0 {
		p = *page
	}

	if limit != nil && *limit > 0 && *limit <= 100 {
		l = *limit
	}

	return PaginationParams{
		Page:  p,
		Limit: l,
	}
}

type PaginatedResult[T any] struct {
	Page       int
	PerPage    int
	PageCount  int
	TotalCount int64
	Data       []T
}
