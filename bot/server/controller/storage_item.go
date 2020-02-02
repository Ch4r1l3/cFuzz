package controller

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StorageItemController struct{}

type UriIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

// swagger:model
type StorageItemExistReq struct {
	// path of storage item in the Image
	// required: true
	// example: /tmp
	Path string `json:"path" binding:"required"`

	// type of storage item
	// required: true
	// example: fuzzer
	Type string `json:"type" binding:"required"`
}

// List StorageItem
func (sic *StorageItemController) List(c *gin.Context) {
	// swagger:operation GET /storage_item storageItem listStorageItem
	// list storageItem
	//
	// list storageItem
	// ---
	// produces:
	// - application/json
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

	storageItems, err := models.GetStorageItems()
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItems)
}

// Retrieve StorageItem
func (sic *StorageItemController) Retrieve(c *gin.Context) {
	// swagger:operation GET /storage_item/{id} storageItem retrieveStorageItem
	// retrieve storageItem
	//
	// retrieve storageItem
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   description: id of StorageItem
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
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	ok, err := models.IsStorageItemExistByID(req.ID)
	if err != nil {
		utils.DBError(c)
	}
	if !ok {
		utils.NotFound(c)
		return
	}
	storageItem, err := models.GetStorageItemByID(req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItem)
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
	//       "$ref": "#/definitions/StorageItemExistReq"
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
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	if !models.IsStorageItemTypeValid(req.Type) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}
	info, err := os.Stat(req.Path)
	if err != nil {
		utils.BadRequestWithMsg(c, "path error, maybe not exists")
		return
	}
	if req.Type == models.Corpus && !info.IsDir() {
		utils.BadRequestWithMsg(c, "corpus should be directory")
		return
	} else if req.Type != models.Corpus && info.IsDir() {
		utils.BadRequestWithMsg(c, "only corpus can be directory")
		return
	}
	if req.Type == models.Target || req.Type == models.Fuzzer {
		os.Chmod(req.Path, 0755)
	}

	storageItem := models.StorageItem{
		Path:          req.Path,
		Type:          req.Type,
		ExistsInImage: true,
	}
	err = models.InsertStorageItem(&storageItem)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, storageItem)
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
	// consumes:
	// - multipart/form-data
	//
	// parameters:
	// - name: type
	//   in: formData
	//   required: true
	//   description: type of storageItem
	//   type: string
	//   enum: [fuzzer, corpus, target]
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

	var Err error

	// get form field
	mtype := c.PostForm("type")
	if mtype == "" {
		utils.BadRequestWithMsg(c, "storageItem type empty")
		return
	}
	if !models.IsStorageItemTypeValid(mtype) {
		utils.BadRequestWithMsg(c, "storageItem type is not valid")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequestWithMsg(c, "please upload file")
		return
	}
	isZipFile := false
	if strings.HasSuffix(header.Filename, ".zip") {
		isZipFile = true
	}

	// create tempory ditrecoy to store file
	tmpDir, err := ioutil.TempDir(config.ServerConf.TempPath, "storageItem")
	if err != nil {
		utils.InternalErrorWithMsg(c, "error create temp directory")
		return
	}
	defer func() {
		if Err != nil {
			os.RemoveAll(tmpDir)
		}
	}()
	var tempFile *os.File
	if isZipFile {
		tempFile, err = ioutil.TempFile(tmpDir, "storageItem.*.zip")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	} else {
		tempFile, err = ioutil.TempFile(tmpDir, "storageItem")
		if err != nil {
			Err = err
			utils.InternalErrorWithMsg(c, "error create temp file")
			return
		}
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		Err = err
		tempFile.Close()
		utils.InternalErrorWithMsg(c, "error copy upload file")
		return
	}
	tempFile.Close()

	//if upload file is zip, unzip it
	if isZipFile {
		err = utils.Unzip(tempFile.Name())
		if err != nil {
			Err = err
			utils.BadRequestWithMsg(c, err.Error())
			return
		}
		os.RemoveAll(tempFile.Name())
	}
	var storePath string
	if mtype == models.Corpus {
		storePath = tmpDir
	} else if isZipFile {
		relPath := c.PostForm("relPath")
		if relPath == "" {
			storePath = filepath.Join(tmpDir, mtype)
		} else {
			storePath = filepath.Join(tmpDir, relPath)
			storePath = filepath.Clean(storePath)
			rel, err := filepath.Rel(tmpDir, storePath)
			if err != nil {
				Err = err
				utils.BadRequestWithMsg(c, "error get rel of targetRelPath")
				return
			}
			if strings.Contains(rel, "..") {
				Err = errors.New("relPath wrong")
				utils.BadRequestWithMsg(c, "path can only under this temp directory")
				return
			}
		}
		if !utils.IsPathExists(storePath) {
			Err = errors.New("zip file should include file")
			utils.BadRequestWithMsg(c, "zip file should include file named "+mtype)
			return
		}
	} else {
		storePath = tempFile.Name()
	}

	// if type is target or fuzzer, make it executable
	if mtype == models.Target || mtype == models.Fuzzer {
		os.Chmod(storePath, 0755)
	}

	storageItem := models.StorageItem{
		Dir:  tmpDir,
		Path: storePath,
		Type: mtype,
	}
	err = models.InsertStorageItem(&storageItem)
	if err != nil {
		utils.DBError(c)
		Err = err
		return
	}
	c.JSON(http.StatusOK, storageItem)
}

// Destroy StorageItem
func (sic *StorageItemController) Destroy(c *gin.Context) {
	// swagger:operation DELETE /storage_item/{id} storageItem deleteStorageItem
	// delete storageItem
	//
	// delete storageItem
	// ---
	// produces:
	// - application/json
	//
	// parameters:
	// - name: id
	//   description: id of StorageItem
	//   in: path
	//   required: true
	//   type: integer
	//
	// responses:
	//   '204':
	//      description: delete success
	//   '403':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"
	//   '404':
	//      description: not found
	//   '500':
	//      schema:
	//        "$ref": "#/definitions/ErrResp"

	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	ok, err := models.IsStorageItemExistByID(req.ID)
	if err != nil {
		utils.DBError(c)
	}
	if !ok {
		utils.NotFound(c)
		return
	}
	storageItem, err := models.GetStorageItemByID(req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	if !storageItem.ExistsInImage {
		os.RemoveAll(storageItem.Dir)
	}
	err = models.DeleteStorageItemByID(req.ID)
	if err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}
