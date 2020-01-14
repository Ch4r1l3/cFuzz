package logger

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
	"path/filepath"
)

const (
	LogFileName = "cfuzz-master.log"
)

var Logger hclog.Logger

func Setup() {
	var writer io.Writer
	if config.ServerConf.LogToFile {
		var logFilePath string
		if config.ServerConf.LogFileDir == "" {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			logFilePath = filepath.Join(filepath.Dir(ex), LogFileName)
		} else {
			stat, err := os.Stat(config.ServerConf.LogFileDir)
			//if log file dir not exist, try to create it
			if os.IsNotExist(err) {
				//if cannot create it
				if err = os.MkdirAll(config.ServerConf.LogFileDir, os.ModePerm); err != nil {
					panic(err)
				}
			} else if !stat.IsDir() {
				panic(errors.New("logfile dir is not a directory"))
			}
			logFilePath = filepath.Join(config.ServerConf.LogFileDir, LogFileName)
		}
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = os.Stdout
	}
	Logger = hclog.New(&hclog.LoggerOptions{
		Name:   "cfuzz-master",
		Output: writer,
		Level:  hclog.Debug,
	})

}
