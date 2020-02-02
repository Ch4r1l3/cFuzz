package utils

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func SaveTempFile(c *gin.Context, formName string, tempFilePrefix string) (string, error) {
	file, header, err := c.Request.FormFile(formName)
	if err != nil {
		BadRequestWithMsg(c, "please upload file")
		return "", err
	}
	prefix := tempFilePrefix
	if strings.HasSuffix(header.Filename, ".zip") {
		prefix += ".*.zip"
	}
	tempFile, err := ioutil.TempFile(config.ServerConf.TempPath, prefix)
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
