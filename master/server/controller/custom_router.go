package controller

import (
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func TaskGetHandler(c *gin.Context) {
	p1 := c.Param("path1")
	p2 := c.Param("path2")
	if p1 != "" && p2 == "" {
		task := new(TaskController)
		if p1 == "count" {
			task.Count(c)
		} else {
			n, err := strconv.ParseUint(p1, 10, 64)
			if err != nil {
				utils.BadRequest(c)
				return
			}
			task.Retrieve(c, n)
		}
	} else if p1 != "" && p2 != "" {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		if p2 == "crash" {
			taskCrash := new(TaskCrashController)
			taskCrash.ListByTaskID(c, n)
		} else if p2 == "result" {
			taskResult := new(TaskResultController)
			taskResult.Retrieve(c, n)
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.NotFound(c)
	}
}

func DeploymentGetHandler(c *gin.Context) {
	p1 := c.Param("path1")
	deploymentController := new(DeploymentController)
	if p1 == "simplist" {
		deploymentController.SimpList(c)
	} else if p1 == "count" {
		deploymentController.Count(c)
	} else {
		n, err := strconv.ParseUint(p1, 10, 64)
		if err != nil {
			utils.BadRequest(c)
			return
		}
		deploymentController.Retrieve(c, n)

	}
}
