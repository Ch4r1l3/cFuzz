package models

type TaskCrash struct {
	ID            uint64 `gorm:"primary_key" json:"id"`
	TaskID        uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	BotCrashID    uint64 `json:"-"`
	Path          string `json:"-"`
	ReproduceAble bool   `json:"reproduceAble"`
}
