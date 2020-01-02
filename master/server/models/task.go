package models

type Task struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	DockerfileID uint64 `json:"dockerfileid"`
	Time         uint64 `json:"time"`
	FuzzerID     uint64 `json:"fuzzerid"`
	Running      bool   `json:"running"`
}

type TaskTarget struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Path   string `json:"-"`
}

type TaskCorpus struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Path   string `json:"-"`
}

type TaskEnvironment struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Value  string `json:"value"`
}

func DeleteTask(taskid uint64) error {
	var err error
	if err = DeleteObjectsByTaskID(&TaskEnvironment{}, taskid); err != nil {
		return err
	}
	if err = DeleteObjectsByTaskID(&TaskArgument{}, taskid); err != nil {
		return err
	}
	if err = DeleteObjectByID(&Task{}, taskid); err != nil {
		return err
	}
	return nil
}

func InsertEnvironments(taskid uint64, environments []string) error {
	for _, v := range environments {
		taskEnvironment := TaskEnvironment{
			TaskID: taskid,
			Value:  v,
		}
		if err := DB.Create(&taskEnvironment).Error; err != nil {
			return err
		}
	}
	return nil
}

type TaskArgument struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func InsertArguments(taskid uint64, arguments map[string]string) error {
	for k, v := range arguments {
		taskArgument := TaskArgument{
			TaskID: taskid,
			Key:    k,
			Value:  v,
		}
		if err := DB.Create(&taskArgument).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Delete(obj).Error
}

type TaskCrash struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
}

type TaskFuzzResult struct {
	Command      string `json:"command"`
	TimeExecuted int    `json:"timeExecuted"`
	TaskID       int    `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
}

type TaskFuzzResultStat struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
}
