package models

type Task struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	DockerfileID uint64 `json:"dockerfileid"`
	Dockerfile   Dockerfile
	Time         uint64 `json:"time"`
	FuzzerID     uint64 `json:"fuzzerid"`
	Fuzzer       Fuzzer
}

type TaskTarget struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid"`
	Task   Task
	Path   string `json:"-"`
}

type TaskCorpus struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid"`
	Task   Task
	Path   string `json:"-"`
}
