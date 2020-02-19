package controller

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type StorageItemController struct{}

type StorageItemTypeReq struct {
	Type string `uri:"type" binding:"required"`
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

func getStorageItem(c *gin.Context) (*models.StorageItem, error) {
	var storageItem models.StorageItem
	err := getObject(c, &storageItem)
	if err != nil {
		return nil, err
	}
	if storageItem.UserID != uint64(c.GetInt64("id")) && !c.GetBool("isAdmin") {
		utils.Forbidden(c)
		return nil, errors.New("no permission")
	}
	return &storageItem, nil
}

// swagger:model
type StorageItemCombine struct {
	Data []models.StorageItem `json:"data"`
	CountResp
}

// List StorageItems
func (sic *StorageItemController) List(c *gin.Context) {
	// swagger:operation GET /storage_item storageItem listStorageItem
	// list storageItem
	//
	// list storageItem
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: offset
	//   in: query
	//   type: integer
	// - name: limit
	//   in: query
	//   type: integer
	// - name: name
	//   in: query
	//   type: string
	//
	// responses:
	//   '200':
	//      schema:
	//        type: array
	//        items:
	//          "$ref": "#/definitions/StorageItem"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var storageItems []models.StorageItem
	count, err := getList(c, &storageItems)
	if err != nil {
		return
	}
	for i, _ := range storageItems {
		if !storageItems[i].ExistsInImage {
			storageItems[i].Path = ""
		}
	}
	c.JSON(http.StatusOK, StorageItemCombine{
		Data: storageItems,
		CountResp: CountResp{
			Count: count,
		},
	})
}

// List StorageItem By Type
func (sic *StorageItemController) ListByType(c *gin.Context) {
	// swagger:operation GET /storage_item/{type} storageItem listStorageItemByType
	// list storageItem by type
	//
	// list storageItem by type
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: type
	//   in: path
	//   required: true
	//   type: string
	// - name: offset
	//   in: query
	//   type: integer
	// - name: limit
	//   in: query
	//   type: integer
	// - name: name
	//   in: query
	//   type: integer
	//
	// responses:
	//   '200':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	mtype := c.Param("type")
	if !models.IsStorageItemTypeValid(mtype) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	offset := c.GetInt("offset")
	limit := c.GetInt("limit")
	name := c.Query("name")
	userID := uint64(c.GetInt64("id"))
	isAdmin := c.GetBool("isAdmin")
	storageItems, count, err := service.GetStorageItemsByTypeCombine(mtype, offset, limit, name, userID, isAdmin)
	for i, _ := range storageItems {
		if !storageItems[i].ExistsInImage {
			storageItems[i].Path = ""
		}
	}
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, StorageItemCombine{
		Data: storageItems,
		CountResp: CountResp{
			Count: count,
		},
	})
}

// Create Exist StorageItem
func (sic *StorageItemController) CreateExist(c *gin.Context) {
	// swagger:operation POST /storage_item/exist storageItem createExistStorageItem
	// create exist storageItem
	//
	// create exist storageItem
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
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req StorageItemExistReq
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	if !models.IsStorageItemTypeValid(req.Type) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	if service.IsStorageItemExistsCombine(req.Name, req.Type, uint64(c.GetInt64("id"))) {
		utils.BadRequestWithMsg(c, "storageItem name exists")
		return
	}
	storageItem := models.StorageItem{
		Name:          req.Name,
		Path:          req.Path,
		Type:          req.Type,
		UserID:        uint64(c.GetInt64("id")),
		ExistsInImage: true,
	}
	err := service.InsertObject(&storageItem)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusCreated, storageItem)
}

// Create StorageItem
func (sic *StorageItemController) Create(c *gin.Context) {
	// swagger:operation POST /storage_item storageItem createStorageItem
	// create storageItem
	//
	// create storageItem
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
	// - name: relPath
	//   description: if upload file is zip and type is not corpus, this field specefiy the path of file like target
	//   in: formData
	//   required: false
	//   type: string
	//
	// responses:
	//   '201':
	//      schema:
	//        "$ref": "#/definitions/StorageItem"
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
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
	if service.IsStorageItemExistsCombine(name, mtype, uint64(c.GetInt64("id"))) {
		utils.BadRequestWithMsg(c, "storageItem name exists")
		return
	}
	var tempFile string
	if tempFile, err = SaveTempFile(c, "file", "storageItem"); err != nil {
		return
	}
	storageItem := models.StorageItem{
		Name:    name,
		Path:    tempFile,
		Type:    mtype,
		RelPath: c.PostForm("relPath"),
		UserID:  uint64(c.GetInt64("id")),
	}
	err = service.InsertObject(&storageItem)
	if err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusCreated, storageItem)
}

// Delete StorageItem
func (sic *StorageItemController) Destroy(c *gin.Context) {
	// swagger:operation Delete /storage_item/{id} storageItem deleteStorageItem
	// delete storageItem
	//
	// delete storageItem
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
	//   '400':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	storageItem, err := getStorageItem(c)
	if err != nil {
		return
	}
	if service.IsStorageItemReferred(storageItem.ID) {
		utils.BadRequestWithMsg(c, "storageItem is being used by task")
		return
	}
	err = service.DeleteStorageItemByID(storageItem.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.String(http.StatusNoContent, "")
}
