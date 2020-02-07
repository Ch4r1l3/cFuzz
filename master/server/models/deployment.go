package models

// swagger:model
type Deployment struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: test-image
	Name string `json:"name"`

	// example: 123
	Content string `json:"content"`

	// example: 1
	UserID uint64 `json:"userID"`
}
