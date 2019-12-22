package controller

import (
	"github.com/gin-gonic/gin"
)

type TaskController struct{}

func (tc *TaskController) Retrieve(c *gin.Context) {

}

func (tc *TaskController) Create(c *gin.Context) {

}

func (tc *TaskController) Destroy(c *gin.Context) {

}

type TaskCrashController struct{}

func (tcc *TaskCrashController) List(c *gin.Context) {

}

type TaskCorpusController struct{}

func (t *TaskCorpusController) Create(c *gin.Context) {

}

func (t *TaskCorpusController) Destroy(c *gin.Context) {

}
