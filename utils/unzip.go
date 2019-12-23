package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func Unzip(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return errors.New("error open zip file")
	}
	dir := filepath.Dir(path)
	for _, file := range reader.File {
		tmpPath := filepath.Join(dir, file.Name)
		tmpPath = filepath.Clean(tmpPath)
		relPath, err := filepath.Rel(dir, tmpPath)
		if err != nil {
			return errors.New("error get rel of zip file")
		}
		if filepath.Join(dir, relPath) != tmpPath {
			return errors.New("zipslip not work here")
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(tmpPath, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return errors.New("error open " + file.Name + " of the zip file")
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return errors.New("error create file when unzip")
		}
		defer targetFile.Close()

		if _, err = io.Copy(targetFile, fileReader); err != nil {
			return errors.New("error copy file")
		}
	}
	return nil
}
