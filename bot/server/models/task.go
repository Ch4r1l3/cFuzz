package models

// swagger:model
type Task struct {
	// example: TaskCreated
	Status string `json:"status"`
	// example: 1
	FuzzerID uint64 `json:"fuzzerID"`
	// example: 2
	CorpusID uint64 `json:"corpusID"`
	// example: 3
	TargetID uint64 `json:"targetID"`
	// example: 60
	FuzzCycleTime uint64 `json:"fuzzCycleTime"` //the fuzz cycle time
	// example: 3600
	MaxTime int `json:"maxTime"` //the total time it runs
}

//Task Status
const (
	TaskCreated = "TaskCreated"
	TaskRunning = "TaskRunning"
	TaskStopped = "TaskStopped"
	TaskError   = "TaskError"
)

func GetTask() (*Task, error) {
	var task Task
	if err := DB.First(&task).Error; err != nil {
		return &Task{}, err
	}
	return &task, nil
}

type TaskArgument struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func InsertArguments(arguments map[string]string) error {
	for k, v := range arguments {
		ta := TaskArgument{
			Key:   k,
			Value: v,
		}
		result := DB.Create(&ta)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetArguments() (map[string]string, error) {
	arguments := make(map[string]string)
	taskArguments := []TaskArgument{}
	if err := DB.Find(&taskArguments).Error; err != nil {
		return nil, err
	}
	for _, v := range taskArguments {
		arguments[v.Key] = v.Value
	}
	return arguments, nil
}

type TaskEnvironment struct {
	Value string `json:"value"`
}

func InsertEnvironments(environments []string) error {
	for _, v := range environments {
		te := TaskEnvironment{
			Value: v,
		}
		result := DB.Create(&te)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetEnvironments() ([]string, error) {
	environments := []string{}
	taskEnvironments := []TaskEnvironment{}
	if err := DB.Find(&taskEnvironments).Error; err != nil {
		return environments, err
	}
	for _, v := range taskEnvironments {
		environments = append(environments, v.Value)
	}
	return environments, nil
}
