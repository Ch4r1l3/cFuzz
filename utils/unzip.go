package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return errors.New("error open zip file")
	}
	defer reader.Close()
	dir := filepath.Dir(path)
	for _, file := range reader.File {
		err = func() error {
			tmpPath := filepath.Join(dir, file.Name)
			tmpPath = filepath.Clean(tmpPath)
			relPath, err := filepath.Rel(dir, tmpPath)
			if err != nil {
				return errors.New("error get rel of zip file")
			}
			if strings.Contains(relPath, "..") {
				return errors.New("zipslip not work here")
			}
			if file.FileInfo().IsDir() {
				os.MkdirAll(tmpPath, file.Mode())
				return nil
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
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
