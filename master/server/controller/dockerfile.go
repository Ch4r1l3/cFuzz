package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DockerfileReq struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content"`
}

type DockerfileUriReq struct {
	Id uint64 `uri:"id" binding:"required"`
}

type DockerfileController struct{}

func (dc *DockerfileController) List(c *gin.Context) {
	dockerfiles := []models.Dockerfile{}
	err := models.GetObjects(&dockerfiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, dockerfiles)
}

func (dc *DockerfileController) Create(c *gin.Context) {
	var req DockerfileReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	err = models.InsertObject(&models.Dockerfile{
		Name:    req.Name,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, "")
}

func (dc *DockerfileController) Update(c *gin.Context) {
	var uriReq DockerfileUriReq
	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	var req DockerfileReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	var dockerfile models.Dockerfile
	err = models.GetObjectById(&dockerfile, uriReq.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	dockerfile.Name = req.Name
	dockerfile.Content = req.Content
	if err = models.DB.Save(&dockerfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusOK, "")
}

func (dc *DockerfileController) Destroy(c *gin.Context) {
	var uriReq DockerfileUriReq

	err := c.ShouldBindUri(&uriReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}
	err = models.DeleteObjectById(models.Dockerfile{}, uriReq.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}
	c.JSON(http.StatusNoContent, "")
}
