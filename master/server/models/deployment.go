package models

type Deployment struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
