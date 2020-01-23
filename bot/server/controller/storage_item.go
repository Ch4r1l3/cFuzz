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

type StorageItemTypeReq struct {
	Type string `json:"type" binding:"required"`
}

type UriIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

type StorageItemExistReq struct {
	Path string `json:"path" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (sic *StorageItemController) List(c *gin.Context) {
	storageItems, err := models.GetStorageItems()
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

func (sic *StorageItemController) Retrieve(c *gin.Context) {
	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequest(c)
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

func (sic *StorageItemController) Create(c *gin.Context) {
	var Err error
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

func (sic *StorageItemController) Destroy(c *gin.Context) {
	var req UriIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		utils.BadRequest(c)
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
