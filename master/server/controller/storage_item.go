package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type StorageItemController struct{}

type StorageItemTypeReq struct {
	Type string `json:"type" binding:"required"`
}

type StorageItemExistReq struct {
	Name string `json:"name" binding:"required"`
	Path string `json:"path" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (sic *StorageItemController) List(c *gin.Context) {
	var storageItems []models.StorageItem
	err := models.GetObjects(&storageItems)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItems)
}

func (sic *StorageItemController) ListByType(c *gin.Context) {
	var req StorageItemTypeReq
	if err := c.ShouldBindUri(&req); err != nil {
		utils.BadRequest(c)
		return
	}
	if !models.IsStorageItemTypeValid(req.Type) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	storageItems, err := models.GetStorageItemsByType(req.Type)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItems)

}

func (sic *StorageItemController) CreateExist(c *gin.Context) {
	var req StorageItemExistReq
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequest(c)
		return
	}
	if !models.IsStorageItemTypeValid(req.Type) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	storageItem := models.StorageItem{
		Name:          req.Name,
		Path:          req.Path,
		Type:          req.Type,
		ExistsInImage: true,
	}
	err := models.InsertObject(&storageItem)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItem)
}

func (sic *StorageItemController) Create(c *gin.Context) {
	var err error
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequestWithMsg(c, "storageItem name empty")
		return
	}
	mtype := c.PostForm("type")
	if mtype == "" {
		utils.BadRequestWithMsg(c, "storageItem type empty")
		return
	}
	if !models.IsStorageItemTypeValid(mtype) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	if models.IsStorageItemExistsByNameAndType(name, mtype) {
		utils.BadRequestWithMsg(c, "storageItem name exists")
		return
	}
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "storageItem"); err != nil {
		return
	}
	storageItem := models.StorageItem{
		Name: name,
		Path: tempFile,
		Type: mtype,
	}
	err = models.InsertObject(&storageItem)
	if err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItem)
}

func (sic *StorageItemController) Destroy(c *gin.Context) {
	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequest(c)
		return
	}
	if !models.IsObjectExistsByID(&models.StorageItem{}, req.ID) {
		utils.NotFound(c)
		return
	}
	var storageItem models.StorageItem
	if err = models.GetObjectByID(&storageItem, req.ID); err != nil {
		utils.DBError(c)
		return
	}
	if !storageItem.ExistsInImage {
		os.RemoveAll(storageItem.Path)
	}
	err = models.DeleteObjectByID(models.StorageItem{}, req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
