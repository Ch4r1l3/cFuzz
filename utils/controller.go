package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// swagger:model
type ErrResp struct {
	// example: some error
	Error string `json:"error"`
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ErrResp{
		Error: "bad request",
	})
}

func DBError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrResp{
		Error: "db error",
	})
}

func BadRequestWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, ErrResp{
		Error: msg,
	})
}

func InternalErrorWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, ErrResp{
		Error: msg,
	})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrResp{
		Error: "not found",
	})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, ErrResp{
		Error: "permission not enough",
	})
}

func ForbiddenWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, ErrResp{
		Error: msg,
	})
}
