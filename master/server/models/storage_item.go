package models

import (
	"os"
)

// item that record some file info
// swagger:model
type StorageItem struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: afl
	Name string `json:"name" sql:"type:varchar(255) NOT NULL UNIQUE"`

	Path string `json:"path"`

	// example: fuzzer
	Type string `json:"type"`

	// example: true
	ExistsInImage bool `json:"existsInImage"`

	// if upload file is zip and type is not corpus, this field specefiy the path of file like target
	// example: test/target
	RelPath string `json:"relPath"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:integer REFERENCES user(id) ON DELETE CASCADE"`
}

func (s *StorageItem) Delete() error {
	if !s.ExistsInImage {
		os.RemoveAll(s.Path)
	}
	return nil
}

// storage item types
const (
	Fuzzer = "fuzzer"
	Target = "target"
	Corpus = "corpus"
)

var StorageItemTypes = []string{Fuzzer, Target, Corpus}

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

func IsStorageItemExistsCombine(name string, mtype string, userID uint64) bool {
	return IsObjectExistsCustom(&StorageItem{}, []string{"name = ?", "type = ?", "user_id = ?"}, []interface{}{name, mtype, userID})
}

func IsStorageItemReferred(id uint64) bool {
	for _, v := range StorageItemTypes {
		if IsObjectExistsCustom(&Task{}, []string{v + "_id = ?"}, []interface{}{id}) {
			return true
		}
	}
	return false
}

func GetStorageItemsByTypeCombine(mtype string, offset int, limit int, name string, userID uint64, isAdmin bool) ([]StorageItem, int, error) {
	var storageItems []StorageItem
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
	var storageItems []StorageItem
	if err := DB.Where(query, id).Find(&storageItems).Error; err != nil {
		return err
	}
	for _, s := range storageItems {
		s.Delete()
	}
	if err := DB.Where(query, id).Delete(StorageItem{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteStorageItemByID(id uint64) error {
	return DeleteStorageItemCustom("id = ?", id)
}
