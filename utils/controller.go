package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "bad request",
	})
}

func DBError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "db error",
	})
}

func BadRequestWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": msg,
	})
}

func InternalErrorWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": msg,
	})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "not found",
	})
}
