package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskResultController struct{}

func (trc *TaskResultController) Retrieve(c *gin.Context) {
	result, stats, err := models.GetFuzzResult()
	if err != nil {
		c.JSON(http.StatusOK, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"command":      result.Command,
		"timeExecuted": result.TimeExecuted,
		"updateAt":     result.UpdateAt,
		"stats":        stats,
	})
}
