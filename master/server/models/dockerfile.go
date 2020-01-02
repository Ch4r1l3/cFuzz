package models

type Dockerfile struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func IsDockerfileExistsByID(id uint64) bool {
	var dockerfiles []Dockerfile
	if err := DB.Where("id = ?", id).Find(&dockerfiles).Error; err != nil {
		return true
	}
	return len(dockerfiles) >= 1
}
