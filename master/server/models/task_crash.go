package models

// swagger:model
type TaskCrash struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: 1
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`

	// example: 1
	BotCrashID uint64 `json:"-"`

	// example: ./crashes/xxxx
	Path string `json:"-"`

	// example: true
	ReproduceAble bool `json:"reproduceAble"`
}
