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
		c.Next()
		return
	}
	c.Set("offset", offset)
	c.Set("limit", limit)
	c.Set("pagination", true)
	c.Next()
}
