package models

// swagger:model
type Image struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: test-image
	Name string `json:"name"`

	// example: true
	IsDeployment bool `json:"isDeployment"`

	// example: 123
	Content string `json:"content"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:bigint REFERENCES user(id) ON DELETE CASCADE"`
}

func IsImageReferred(id uint64) bool {
	return IsObjectExistsCustom(&Task{}, []string{"image_id = ?"}, []interface{}{id})
}
