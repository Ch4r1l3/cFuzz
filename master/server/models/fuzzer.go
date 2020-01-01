package models

type Fuzzer struct {
	ID   uint64 `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	Path string `json:"-"`
}

func IsFuzzerExists(name string) bool {
	var fuzzer []Fuzzer
	if err := DB.Where("name = ?", name).Find(&fuzzer).Error; err != nil {
		return true
	}
	return len(fuzzer) >= 1
}
