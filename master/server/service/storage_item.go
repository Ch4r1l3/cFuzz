package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
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
		count, err = GetObjectCombinCustom(&storageItems, offset, limit, name, []string{"type = ?"}, []interface{}{mtype})
	} else {
		count, err = GetObjectCombinCustom(&storageItems, offset, limit, name, []string{"type = ?", "user_id = ?"}, []interface{}{mtype, userID})
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
