package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TaskDeleteHandler(c *gin.Context) {
	p1 := c.Param("path1")
	p2 := c.Param("path2")
	p3 := c.Param("path3")
	if p1 != "" && p2 == "" && p3 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}
		task := new(TaskController)
		task.Destroy(c, n)
	} else if p1 != "" && p2 == "target" && p3 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}

		taskTarget := new(TaskTargetController)
		taskTarget.Destroy(c, n)
	} else if p1 != "" && p2 == "corpus" && p3 != "" {
		n1, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}

		n2, err := strconv.ParseUint(p3, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}
		taskCorpus := new(TaskCorpusController)
		taskCorpus.Destroy(c, n1, n2)
	}
}

type TaskController struct{}

func (tc *TaskController) List(c *gin.Context) {

}

func (tc *TaskController) Create(c *gin.Context) {

}

func (tc *TaskController) Update(c *gin.Context) {

}

func (tc *TaskController) Destroy(c *gin.Context, id uint64) {

}

type TaskCorpusController struct{}

func (tcc *TaskCorpusController) List(c *gin.Context) {

}

func (tcc *TaskCorpusController) Create(c *gin.Context) {

}

func (tcc *TaskCorpusController) Destroy(c *gin.Context, taskid uint64, corpusid uint64) {

}

type TaskTargetController struct{}

func (ttc *TaskTargetController) Retrieve(c *gin.Context) {

}

func (ttc *TaskTargetController) Create(c *gin.Context) {

}

func (ttc *TaskTargetController) Destroy(c *gin.Context, id uint64) {

}
