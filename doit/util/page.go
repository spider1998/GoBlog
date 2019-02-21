package util

import (
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
)

const (
	DefaultPageSize int = 100
	MaxPageSize     int = 1000
)

type pager struct {
	Page       int
	PageSize   int
	TotalCount int
}

func (p *pager) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *pager) Limit() int {
	return p.PageSize
}

func newPager(page, perPage, total int) *pager {
	if perPage < 1 {
		perPage = 100
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &pager{
		Page:       page,
		PageSize:   perPage,
		TotalCount: total,
	}
}

func GetPaginatedListFromRequest(c *routing.Context, count int) *pager {
	page := ParseInt(c.Query("page"), 1)
	pageSize := ParseInt(c.Query("page_size"), DefaultPageSize)
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	p := newPager(page, pageSize, count)

	c.Response.Header().Set("X-Page-Total", strconv.Itoa(p.TotalCount))
	c.Response.Header().Set("X-Page", strconv.Itoa(p.Page))
	c.Response.Header().Set("X-Page-Size", strconv.Itoa(p.PageSize))

	return p
}

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}


func ParsePagination(c *routing.Context) (page, pageSize int64, err error) {
	page = 1
	pageSize = 50
	if pageStr := c.Query("page"); pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			return
		}
	}
	return
}