package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/jinzhu/gorm"
)

func IsStorageItemExistsCombine(name string, mtype string, userID uint64) bool {
	return IsObjectExistsCustom(&models.StorageItem{}, []string{"name = ?", "type = ?", "user_id = ?"}, []interface{}{name, mtype, userID})
}

func IsStorageItemReferred(id uint64) bool {
	for _, v := range models.StorageItemTypes {
		if IsObjectExistsCustom(&models.Task{}, []string{v + "_id = ?"}, []interface{}{id}) {
			return true
		}
	}
	return false
}

func GetStorageItemsByTypeCombine(mtype string, offset int, limit int, name string, userID uint64, isAdmin bool) ([]models.StorageItem, int, error) {
	var storageItems []models.StorageItem
	var count int
	var err error
	if isAdmin {
		count, err = getObjectCombinCustom(&storageItems, offset, limit, name, []string{"type = ?"}, []interface{}{mtype})
	} else {
		count, err = getObjectCombinCustom(&storageItems, offset, limit, name, []string{"type = ?", "user_id = ?"}, []interface{}{mtype, userID})
	}
	return storageItems, count, err
}

func DeleteStorageItemCustom(query string, id uint64) error {
	var storageItems []models.StorageItem
	if err := models.DB.Where(query, id).Find(&storageItems).Error; err != nil {
		return err
	}
	for _, s := range storageItems {
		s.Delete()
	}
	if err := models.DB.Where(query, id).Delete(models.StorageItem{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteStorageItemByID(id uint64) error {
	return DeleteStorageItemCustom("id = ?", id)
}

func CreateStorageItem(storageItem *models.StorageItem) error {
	return insertObject(storageItem)
}

func GetStorageItemByID(id uint64) (*models.StorageItem, error) {
	var storageItem models.StorageItem
	if err := GetObjectByID(&storageItem, id); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &storageItem, nil
}
