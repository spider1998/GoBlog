package util

import (
	"github.com/go-ozzo/ozzo-routing"
	"strconv"
)

//分页处理
func ParsePagination(c *routing.Context) (page, pageSize int64, err error) {
	page = 1
	pageSize = 50
	//获取具体页数
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
