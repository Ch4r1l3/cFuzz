package utils

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
)

func SaveTempFile(c *gin.Context, formName string, tempFilePrefix string) (string, error) {
	file, _, err := c.Request.FormFile(formName)
	if err != nil {
		BadRequest(c)
		return "", err
	}
	tempFile, err := ioutil.TempFile(config.ServerConf.TempPath, tempFilePrefix)
	if err != nil {
		InternalErrorWithMsg(c, "error create temp file")
		return "", err
	}
	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()
		os.RemoveAll(tempFile.Name())
		InternalErrorWithMsg(c, "error copy upload file")
		return "", err
	}
	tempFile.Close()
	return tempFile.Name(), nil
}
