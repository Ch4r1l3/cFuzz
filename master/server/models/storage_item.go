package models

type StorageItem struct {
	ID            uint64 `gorm:"primary_key" json:"id"`
	Name          string `json:"name"`
	Path          string `json:"-"`
	Type          string `json:"type"`
	ExistsInImage bool   `json:"existsInImage"`
}

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
	if err := DB.Where("type = ?", mtype).Find(&storageItems).Error; err != nil {
		return nil, err
	}
	return storageItems, nil
}
