package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeploymentReq struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content"`
}

type DeploymentUriReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

type DeploymentController struct{}

func (dc *DeploymentController) List(c *gin.Context) {
	deployments := []models.Deployment{}
	err := models.GetObjects(&deployments)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, deployments)
}

func (dc *DeploymentController) Create(c *gin.Context) {
	var req DeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	deployment := models.Deployment{
		Name:    req.Name,
		Content: req.Content,
	}
	err = models.InsertObject(&deployment)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, deployment)
}

func (dc *DeploymentController) Update(c *gin.Context) {
	var uriReq DeploymentUriReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var req DeploymentReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var deployment models.Deployment
	err = models.GetObjectByID(&deployment, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	deployment.Name = req.Name
	deployment.Content = req.Content
	if err = models.DB.Save(&deployment).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, "")
}

func (dc *DeploymentController) Destroy(c *gin.Context) {
	var uriReq DeploymentUriReq

	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	err = models.DeleteObjectByID(models.Deployment{}, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
