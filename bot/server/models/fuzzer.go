package models

type Fuzzer struct {
	ID   uint64 `gorm:"primary_key" json:"id"`
	Path string `json:"-"`
	Name string `json:"name"`
}

func GetFuzzerByID(id uint64) (*Fuzzer, error) {
	var fuzzer Fuzzer
	if err := DB.Where("id = ?", id).First(&fuzzer).Error; err != nil {
		return &Fuzzer{}, err
	}
	return &fuzzer, nil
}
