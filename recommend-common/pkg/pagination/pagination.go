package pagination

import "math"

type Pageable struct {
	PageNumber *uint64 `json:"pageNumber" form:"pageNumber"`
	PageSize   *uint64 `json:"pageSize"   form:"pageSize"`
}

func (pageable *Pageable) GetPageSize() uint64 {
	if pageable.PageSize == nil || *pageable.PageSize <= 0 {
		return 10
	}
	return *pageable.PageSize
}

func (pageable *Pageable) GetPageNumber() uint64 {
	if pageable.PageNumber == nil || *pageable.PageNumber <= 1 {
		return 1
	}
	return *pageable.PageNumber
}

func (pageable *Pageable) GetOffset() uint64 {
	return pageable.GetPageSize() * (pageable.GetPageNumber() - 1)
}

type Paged[T any] struct {
	Content       []*T   `json:"content"`
	PageSize      uint64 `json:"pageSize"`
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func NewPaged[T any](content []*T, pageSize uint64, totalElements uint64) Paged[T] {
	return Paged[T]{
		Content:       content,
		PageSize:      pageSize,
		TotalPages:    uint64(math.Ceil(float64(totalElements) / float64(pageSize))),
		TotalElements: totalElements,
	}
}
