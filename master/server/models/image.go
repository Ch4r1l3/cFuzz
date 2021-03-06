package models

// swagger:model
type Image struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: test-image
	Name string `json:"name" sql:"type:varchar(255) NOT NULL UNIQUE"`

	// example: true
	IsDeployment bool `json:"isDeployment"`

	// example: 123
	Content string `json:"content"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:integer REFERENCES user(id) ON DELETE CASCADE"`
}
