package models

type Dockerfile struct {
	ID      uint64 `gorm:"primary_key";json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func InsertDockerfile(name string, content string) error {
	dockerfile := Dockerfile{
		Name:    name,
		Content: content,
	}
	if err := DB.Create(&dockerfile).Error; err != nil {
		return err
	}
	return nil
}

func GetDockerfiles() ([]Dockerfile, error) {
	dockerfiles := []Dockerfile{}
	if err := DB.Find(&dockerfiles).Error; err != nil {
		return nil, err
	}
	return dockerfiles, nil
}

func UpdateDockerfile(id uint64, name string, content string) error {
	var dockerfile Dockerfile
	if err := DB.Where("id = ?", id).First(&dockerfile).Error; err != nil {
		return err
	}
	if name != "" {
		if err := DB.Model(dockerfile).Update("Name", name).Error; err != nil {
			return err
		}
	}
	if content != "" {
		if err := DB.Model(dockerfile).Update("Content", content).Error; err != nil {

			return err
		}
	}
	return nil
}
