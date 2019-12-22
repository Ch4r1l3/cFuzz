package models

import "github.com/jinzhu/gorm"

type Task struct {
	CorpusDir  string `json:"corpusDir"`
	TargetPath string `json:"targetPath"`
}

type TaskArguments struct {
	gorm.Model
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TaskEnviroments struct {
	gorm.Model
	Value string `json:"value"`
}
