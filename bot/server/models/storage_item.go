package models

import (
	"github.com/jinzhu/gorm"
)

// StorageItem ...
// swagger:model
type StorageItem struct {
	// in: body
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`
	// example: /tmp/storageItem123x/
	Dir string `json:"-"` // Directory of StorageItem
	// example: /tmp/storageItem123x/123
	Path string `json:"-"` // Path of StorageItem
	// example: Fuzzer
	Type string `json:"type"`
	// example: true
	ExistsInImage bool `json:"existsInImage"` //whether the StorageItem exist in the image
}

//types of StorageItem
const (
	Fuzzer = "fuzzer"
	Target = "target"
	Corpus = "corpus"
)

// check the input type
func IsStorageItemTypeValid(mtype string) bool {
	switch mtype {
	case
		Fuzzer,
		Target,
		Corpus:
		return true
	}
	return false
}

func GetStorageItemsByType(mtype string) ([]StorageItem, error) {
	var storageItems []StorageItem
	if err := DB.Where("type = ?", mtype).Find(&storageItems).Error; err != nil {
		return nil, err
	}
	return storageItems, nil
}

func IsStorageItemExistByID(id uint64) (bool, error) {
	var storageItem StorageItem
	if err := DB.Where("id = ?", id).First(&storageItem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetStorageItemByID(id uint64) (*StorageItem, error) {
	var storageItem StorageItem
	if err := DB.Where("id = ?", id).First(&storageItem).Error; err != nil {
		return nil, err
	}
	return &storageItem, nil
}

func GetStorageItems() ([]StorageItem, error) {
	var storageItems []StorageItem
	if err := DB.Find(&storageItems).Error; err != nil {
		return nil, err
	}
	return storageItems, nil
}

func InsertStorageItem(storageItem *StorageItem) error {
	return DB.Create(storageItem).Error
}

func DeleteStorageItemByID(id uint64) error {
	return DB.Where("id = ?", id).Delete(&StorageItem{}).Error
}
