package models

type Fuzzer struct {
	ID   uint64 `gorm:"primary_key";json:"id"`
	Path string `json:"-"`
	Name string `json:"name"`
}

func GetFuzzerByName(name string) (*Fuzzer, error) {
	var fuzzer Fuzzer
	if err := DB.Where("name = ?", name).First(&fuzzer).Error; err != nil {
		return &Fuzzer{}, err
	}
	return &fuzzer, nil
}
