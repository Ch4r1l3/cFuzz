package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DockerfileReq struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content"`
}

type DockerfileUriReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

type DockerfileController struct{}

func (dc *DockerfileController) List(c *gin.Context) {
	dockerfiles := []models.Dockerfile{}
	err := models.GetObjects(&dockerfiles)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, dockerfiles)
}

func (dc *DockerfileController) Create(c *gin.Context) {
	var req DockerfileReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	dockerfile := models.Dockerfile{
		Name:    req.Name,
		Content: req.Content,
	}
	err = models.InsertObject(&dockerfile)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, dockerfile)
}

func (dc *DockerfileController) Update(c *gin.Context) {
	var uriReq DockerfileUriReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var req DockerfileReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	var dockerfile models.Dockerfile
	err = models.GetObjectByID(&dockerfile, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	dockerfile.Name = req.Name
	dockerfile.Content = req.Content
	if err = models.DB.Save(&dockerfile).Error; err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, "")
}

func (dc *DockerfileController) Destroy(c *gin.Context) {
	var uriReq DockerfileUriReq

	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	err = models.DeleteObjectByID(models.Dockerfile{}, uriReq.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
