package util

import (
	"math"
)

// Paginator is an interface for pagination utils
type Paginator interface {
	CreateProperties(totalData int) *PaginationDetails
}

// Pagination struct
type Pagination struct {
	Limit  int
	Offset int
	Sort   map[string]string
	Page   int
	Filter map[string]string
	Search string
}

// PaginationDetails details of pagination struct
type PaginationDetails struct {
	Limit       int   `json:"limit"`
	CurrentPage int   `json:"current_page"`
	TotalData   int   `json:"total_data"`
	TotalPage   int64 `json:"total_page"`
	Offset      int   `json:"-"`
}

// CreateProperties calculate pagination data
func (paging *Pagination) CreateProperties(totalData int) *PaginationDetails {
	var offset int
	var totalPage int64
	var countPage float64
	var currentPage int

	if paging.Limit != 0 {
		offset = ((paging.Page * paging.Limit) - paging.Limit)
		countPage = float64(totalData) / float64(paging.Limit)
		if (countPage) > 1 {
			if (totalData % paging.Limit) > 0 {
				math.Floor(countPage)
				totalPage = int64(countPage) + 1
			} else {
				math.Floor(countPage)
				totalPage = int64(countPage)
			}
		} else {
			totalPage = 1
		}
		currentPage = paging.Page
		if offset < 0 {
			currentPage = 1
		}
	}

	data := PaginationDetails{
		Offset:      offset,
		TotalPage:   totalPage,
		TotalData:   totalData,
		CurrentPage: currentPage,
	}

	return &data
}
