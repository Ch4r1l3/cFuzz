package models

// item that record some file info
// swagger:model
type StorageItem struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: afl
	Name string `json:"name"`

	Path string `json:"path"`

	// example: fuzzer
	Type string `json:"type"`

	// example: true
	ExistsInImage bool `json:"existsInImage"`

	// if upload file is zip and type is not corpus, this field specefiy the path of file like target
	// example: test/target
	RelPath string `json:"relPath"`

	// example: 1
	UserID uint64 `json:"userID"`
}

// storage item types
const (
	Fuzzer = "fuzzer"
	Target = "target"
	Corpus = "corpus"
)

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

func IsStorageItemExistsByNameAndType(name string, mtype string) bool {
	var storageItems []StorageItem
	if err := DB.Where("name = ? AND type = ?", name, mtype).Find(&storageItems).Error; err != nil {
		return true
	}
	return len(storageItems) >= 1
}

func GetStorageItemsByType(mtype string) ([]StorageItem, error) {
	var storageItems []StorageItem
	if err := DB.Where("type = ?", mtype).Order("id").Find(&storageItems).Error; err != nil {
		return nil, err
	}
	return storageItems, nil
}

func GetStorageItemsByTypeCombine(mtype string, offset int, limit int, name string) ([]StorageItem, int, error) {
	var storageItems []StorageItem
	var count int
	t := DB.Order("id")
	if name != "" {
		if err := DB.Model(&storageItems).Where("type = ? AND name LIKE ?", mtype, "%"+name+"%").Count(&count).Error; err != nil {
			return nil, 0, err
		}
		t = t.Where("type = ? AND name LIKE ?", mtype, "%"+name+"%")
	} else {
		if err := DB.Model(&storageItems).Where("type = ?", mtype).Count(&count).Error; err != nil {
			return nil, 0, err
		}
		t = t.Where("type = ?", mtype)
	}
	if offset >= 0 && limit >= 0 {
		t = t.Offset(offset).Limit(limit)
	}
	err := t.Find(&storageItems).Error
	return storageItems, count, err
}

func GetStorageItemsByTypeAndUserIDCombine(mtype string, offset int, limit int, name string, userID uint64) ([]StorageItem, int, error) {
	var storageItems []StorageItem
	var count int
	t := DB.Order("id")
	if name != "" {
		if err := DB.Model(&storageItems).Where("type = ? AND name LIKE ? AND user_id = ?", mtype, "%"+name+"%", userID).Count(&count).Error; err != nil {
			return nil, 0, err
		}
		t = t.Where("type = ? AND name LIKE ? AND user_id = ?", mtype, "%"+name+"%", userID)
	} else {
		if err := DB.Model(&storageItems).Where("type = ? AND user_id = ?", mtype, userID).Count(&count).Error; err != nil {
			return nil, 0, err
		}
		t = t.Where("type = ? AND user_id = ?", mtype, userID)
	}
	if offset >= 0 && limit >= 0 {
		t = t.Offset(offset).Limit(limit)
	}
	err := t.Find(&storageItems).Error
	return storageItems, count, err
}
