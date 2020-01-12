package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

type TaskCorpusController struct{}

func (tcc *TaskCorpusController) Retrieve(c *gin.Context, taskid uint64) {
	var corpus []models.TaskCorpus
	if err := models.GetObjectsByTaskID(&corpus, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, corpus)
}

func (tcc *TaskCorpusController) Create(c *gin.Context) {
	var taskid uint64
	var err error
	if taskid, err = getTaskID(c); err != nil {
		return
	}
	var corpusArray []models.TaskCorpus
	if err = models.GetObjectsByTaskID(&corpusArray, taskid); err != nil {
		utils.DBError(c)
		return
	}
	if len(corpusArray) > 0 {
		utils.BadRequestWithMsg(c, "you should delete corpus first")
		return
	}
	var tempFile string
	if tempFile, err = utils.SaveTempFile(c, "file", "corpus"); err != nil {
		return
	}
	corpus := models.TaskCorpus{
		TaskID:   taskid,
		Path:     tempFile,
		FileName: filepath.Base(tempFile),
	}
	if err := models.InsertObject(&corpus); err != nil {
		os.RemoveAll(tempFile)
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusOK, corpus)
}

func (tcc *TaskCorpusController) Destroy(c *gin.Context, taskid uint64) {
	var corpus []models.TaskCorpus
	if err := models.GetObjectsByTaskID(&corpus, taskid); err != nil {
		utils.DBError(c)
		return
	}
	for _, v := range corpus {
		os.RemoveAll(v.Path)
	}
	if err := models.DeleteObjectsByTaskID(&models.TaskCorpus{}, taskid); err != nil {
		utils.DBError(c)
		return
	}
	c.JSON(http.StatusNoContent, "")
}

/*
func (tcc *TaskCorpusController) DestroyByID(c *gin.Context, taskid uint64, corpusid uint64) {
        if !models.IsObjectExistsByID(&models.Task{}, taskid) {
                utils.NotFound(c)
                return
        }
        if !models.IsObjectExistsByID(&models.TaskCorpus{}, corpusid) {
                utils.NotFound(c)
                return
        }
        if err := models.DeleteObjectByID(&models.TaskCorpus{}, corpusid); err != nil {
                utils.DBError(c)
                return
        }
        c.JSON(http.StatusNoContent, "")
}
*/
