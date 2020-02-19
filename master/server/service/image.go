package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/jinzhu/gorm"
)

func IsImageReferred(id uint64) bool {
	return IsObjectExistsCustom(&models.Task{}, []string{"image_id = ?"}, []interface{}{id})
}

func CreateImage(image *models.Image) error {
	return insertObject(image)
}

func GetImageByID(id uint64) (*models.Image, error) {
	var image models.Image
	if err := GetObjectByID(&image, id); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func UpdateImage(image *models.Image) error {
	return SaveObject(image)
}
