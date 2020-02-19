package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
)

func IsImageReferred(id uint64) bool {
	return IsObjectExistsCustom(&models.Task{}, []string{"image_id = ?"}, []interface{}{id})
}

func CreateImage(image *models.Image) error {
	return InsertObject(image)
}
