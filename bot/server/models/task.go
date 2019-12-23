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
		return nil, err
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

type TaskEnviroment struct {
	Value string `json:"value"`
}

func InsertEnviroments(enviroments []string) error {
	for _, v := range enviroments {
		te := TaskEnviroment{
			Value: v,
		}
		result := DB.Create(&te)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetEnviroments() ([]string, error) {
	enviroments := []string{}
	taskEnviroments := []TaskEnviroment{}
	if err := DB.Find(&taskEnviroments).Error; err != nil {
		return enviroments, err
	}
	for _, v := range taskEnviroments {
		enviroments = append(enviroments, v.Value)
	}
	return enviroments, nil
}

type TaskCrash struct {
	ID   uint64 `gorm:"primary_key";json:"id"`
	Url  string `json:"url"`
	Path string `json:"-"`
}

func GetCrashes() ([]TaskCrash, error) {
	taskCrashes := []TaskCrash{}
	if err := DB.Find(&taskCrashes).Error; err != nil {
		return nil, err
	}
	return taskCrashes, nil
}
