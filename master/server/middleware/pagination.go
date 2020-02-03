package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Pagination(c *gin.Context) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = -1
	}
	c.Set("offset", offset)
	c.Set("limit", limit)
	c.Next()
}
