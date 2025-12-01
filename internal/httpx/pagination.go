package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)


type PaginationParams struct {
	Page int32
	Limit int32
}

type Pagination struct {
	Offset int32    `json:"offset"`
	Limit int32     `json:"limit"`
	CurrentPage int32 `json:"current_page"`
}

func calculatePagination(page, perPage int32) Pagination {
    if page <= 0{
		page = 1 
	}

	if perPage <= 0 {
		perPage = 10
	} else if perPage > 100 {
		perPage = 100
	}

	offset := int32((page - 1) * perPage)
	return Pagination{
		Offset: int32(offset),
		Limit: int32(perPage),
		CurrentPage: int32(page),
	}
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	return PaginationParams{
		Page: int32(page),
		Limit: int32(limit),
	}
}

func GetPaginationFromQuery(c *gin.Context) Pagination {
	params := GetPaginationParams(c)
	return calculatePagination(params.Page, params.Limit)
}


