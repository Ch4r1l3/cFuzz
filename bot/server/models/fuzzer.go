package models

import "github.com/jinzhu/gorm"

type Fuzzer struct {
	gorm.Model
	Path string `json:"-"`
	Name string `json:"name"`
}
