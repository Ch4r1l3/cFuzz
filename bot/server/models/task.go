package models

type Task struct {
	CorpusDir  string `json:"corpusDir"`
	TargetPath string `json:"targetPath"`
}

type TaskArguments struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TaskEnviroments struct {
	Value string `json:"value"`
}
