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

func GetStorageItemsByTypePagination(mtype string, offset int, limit int) ([]StorageItem, error) {
	var storageItems []StorageItem
	if err := DB.Where("type = ?", mtype).Order("id").Offset(offset).Limit(limit).Find(&storageItems).Error; err != nil {
		return nil, err
	}
	return storageItems, nil
}
