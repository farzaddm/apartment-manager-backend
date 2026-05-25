package pagination

import (
	"math"
)

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int64) *PagedList[T] {
	pl := &PagedList[T]{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviousPage = pl.PageNumber > 1

	return pl
}

// Paginate
func Paginate[TInput any, TOutput any](totalRows int64, items *[]TOutput, pageNumber int, pageSize int64) (*PagedList[TOutput], error) {
	return NewPagedList(items, totalRows, pageNumber, pageSize), nil

}

type PagedList[T any] struct {
	PageNumber      int   `json:"pageNumber"`
	PageSize        int64 `json:"pageSize"`
	TotalRows       int64 `json:"totalRows"`
	TotalPages      int   `json:"totalPages"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
	HasNextPage     bool  `json:"hasNextPage"`
	Items           *[]T  `json:"items"`
}
