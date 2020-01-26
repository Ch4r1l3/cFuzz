package models

// swagger:model
type TaskCrash struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`
	// example: /tmp/afl/123
	Path string `json:"path"`
	// example: true
	ReproduceAble bool `json:"reproduceAble"`
}

func GetCrashes() ([]TaskCrash, error) {
	taskCrashes := []TaskCrash{}
	if err := DB.Find(&taskCrashes).Error; err != nil {
		return nil, err
	}
	return taskCrashes, nil
}

func GetCrashByID(id uint64) (*TaskCrash, error) {
	var crash TaskCrash
	if err := DB.Where("id = ?", id).First(&crash).Error; err != nil {
		return nil, err
	}
	return &crash, nil
}

func CreateCrash(path string, reproduceAble bool) error {
	taskCrash := TaskCrash{
		Path:          path,
		ReproduceAble: reproduceAble,
	}
	return DB.Save(&taskCrash).Error
}
