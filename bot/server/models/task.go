package models

type Task struct {
	CorpusDir  string `json:"corpusDir"`
	TargetDir  string `json:"targetDir"`
	TargetPath string `json:"targetPath"`
	Status     string `json:"status"`
	FuzzerName string `json:"fuzzerName"`
	MaxTime    int    `json:"maxTime"`
}

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

type TaskCrash struct {
	ID            uint64 `gorm:"primary_key";json:"id"`
	Url           string `json:"url"`
	Path          string `json:"-"`
	ReproduceAble bool   `json:"reproduceAble"`
}

func GetCrashes() ([]TaskCrash, error) {
	taskCrashes := []TaskCrash{}
	if err := DB.Find(&taskCrashes).Error; err != nil {
		return nil, err
	}
	return taskCrashes, nil
}

func CreateCrash(path string, reproduceAble bool) error {
	taskCrash := TaskCrash{
		Path:          path,
		ReproduceAble: reproduceAble,
	}
	return DB.Save(&taskCrash).Error
}

type TaskFuzzResult struct {
	Command      string `json:"command"`
	TimeExecuted int    `json:"timeExecuted"`
}

type TaskFuzzResultStat struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func CreateFuzzResult(command []string, stats map[string]string, timeExecuted int) error {
	DB.Delete(&TaskFuzzResult{})
	DB.Delete(&TaskFuzzResultStat{})
	tcommand := ""
	for _, v := range command {
		tcommand += v + " "
	}
	fuzzResult := TaskFuzzResult{
		Command:      tcommand,
		TimeExecuted: timeExecuted,
	}
	if err := DB.Save(&fuzzResult).Error; err != nil {
		return err
	}
	for k, v := range stats {
		stat := TaskFuzzResultStat{
			Key:   k,
			Value: v,
		}
		if err := DB.Save(&stat).Error; err != nil {
			return err
		}
	}
	return nil
}
