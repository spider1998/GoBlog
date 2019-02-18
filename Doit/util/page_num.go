package util

import (
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
)

type Pager struct {
	page       int
	pageSize   int
	totalCount int
}

func (p *Pager) Offset() int {
	return (p.page - 1) * p.pageSize
}

func (p *Pager) Limit() int {
	return p.pageSize
}

func newPagers(page, perPage, total int) *Pager {
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

	return &Pager{
		page:       page,
		pageSize:   perPage,
		totalCount: total,
	}
}

func NewPagerFromRequest(c *routing.Context, count int) *Pager {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), DefaultPageSize)
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	p := newPagers(page, pageSize, count)

	c.Response.Header().Set("X-Page-Total", strconv.Itoa(p.totalCount))
	c.Response.Header().Set("X-Page", strconv.Itoa(p.page))
	c.Response.Header().Set("X-Page-Size", strconv.Itoa(p.pageSize))

	return p
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}
