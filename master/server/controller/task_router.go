package controller

import (
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func TaskDeleteHandler(c *gin.Context) {
	p1 := c.Param("path1")
	p2 := c.Param("path2")
	if p1 != "" && p2 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		task := new(TaskController)
		task.Destroy(c, n)
	} else if p1 != "" && p2 != "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		if p2 == "target" {
			taskTarget := new(TaskTargetController)
			taskTarget.Destroy(c, n)
		} else if p2 == "corpus" {
			taskCorpus := new(TaskCorpusController)
			taskCorpus.Destroy(c, n)
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.NotFound(c)
	}
}

func TaskGetHandler(c *gin.Context) {
	p1 := c.Param("path1")
	p2 := c.Param("path2")
	if p1 != "" && p2 == "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		task := new(TaskController)
		task.Retrieve(c, n)
	} else if p1 != "" && p2 != "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		if p2 == "target" {
			taskTarget := new(TaskTargetController)
			taskTarget.Retrieve(c, n)
		} else if p2 == "corpus" {
			taskCorpus := new(TaskCorpusController)
			taskCorpus.Retrieve(c, n)
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.NotFound(c)
	}
}
