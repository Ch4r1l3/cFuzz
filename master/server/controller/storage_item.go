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

// swagger:model
type StorageItemExistReq struct {
	// example: afl
	Name string `json:"name" binding:"required"`
	// example: /tmp/afl/123
	Path string `json:"path" binding:"required"`
	// example: fuzzer
	Type string `json:"type" binding:"required"`
}

// List StorageItems
func (sic *StorageItemController) List(c *gin.Context) {
	// swagger:operation GET /storage_item storageItem listStorageItem
	// List StorageItem
	//
	// ---
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var storageItems []models.StorageItem
	err := models.GetObjects(&storageItems)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItems)
}

// List StorageItem By Type
func (sic *StorageItemController) ListByType(c *gin.Context) {
	// swagger:operation GET /storage_item/{type} storageItem listStorageItemByType
	// List StorageItem By Type
	//
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: type
	//   in: path
	//   required: true
	//   type: string
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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

// Create Exist StorageItem
func (sic *StorageItemController) CreateExist(c *gin.Context) {
	// swagger:operation POST /storage_item/exist storageItem createExistStorageItem
	// Create Exist StorageItem
	//
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: storageItemExistReq
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/StorageItemExistReq"
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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

// Create StorageItem
func (sic *StorageItemController) Create(c *gin.Context) {
	// swagger:operation POST /storage_item storageItem createStorageItem
	// Create StorageItem
	//
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: name
	//   in: formData
	//   required: true
	//   type: string
	// - name: type
	//   in: formData
	//   required: true
	//   type: string
	// - name: file
	//   in: formData
	//   required: true
	//   type: file
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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

// Delete StorageItem
func (sic *StorageItemController) Destroy(c *gin.Context) {
	// swagger:operation Delete /storage_item/{id} storageItem deleteStorageItem
	// Delete StorageItem
	//
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

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
