package models

type Task struct {
	ID            uint64 `gorm:"primary_key" json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	DeploymentID  uint64 `json:"deploymentid"`
	Time          uint64 `json:"time"`
	FuzzCycleTime uint64 `json:"fuzzCycleTime"`
	FuzzerID      uint64 `json:"fuzzerid"`
	Running       bool   `json:"running"`
}

type TaskTarget struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Path   string `json:"-"`
}

type TaskCorpus struct {
	ID       uint64 `gorm:"primary_key" json:"id"`
	TaskID   uint64 `json:"taskid" sql:"type:bigint REFERENCES task(id) ON DELETE CASCADE"`
	Path     string `json:"-"`
	FileName string `json:"filename"`
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

func GetEnvironments(taskid uint64) ([]string, error) {
	var taskEnvironments []TaskEnvironment
	if err := DB.Where("task_id = ?", taskid).Find(&taskEnvironments).Error; err != nil {
		return nil, err
	}
	environments := []string{}
	for _, v := range taskEnvironments {
		environments = append(environments, v.Value)
	}
	return environments, nil
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

func GetArguments(taskid uint64) (map[string]string, error) {
	var taskArguments []TaskArgument
	if err := DB.Where("task_id = ?", taskid).Find(&taskArguments).Error; err != nil {
		return nil, err
	}
	arguments := make(map[string]string)
	for _, v := range taskArguments {
		arguments[v.Key] = v.Value
	}
	return arguments, nil
}

func DeleteObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Delete(obj).Error
}

func GetObjectsByTaskID(obj interface{}, taskid uint64) error {
	return DB.Where("task_id = ?", taskid).Find(obj).Error
}

func IsObjectExistsByTaskID(obj interface{}, taskid uint64) bool {
	if err := DB.Where("task_id = ?", taskid).First(obj).Error; err != nil {
		return false
	}
	return true
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
