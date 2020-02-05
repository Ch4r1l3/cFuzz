package logger

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
	"path/filepath"
)

const (
	LogFileName = "cfuzz-bot.log"
)

var Logger hclog.Logger

func Setup() {
	var writer io.Writer
	var logFilePath string
	stat, err := os.Stat(config.ServerConf.TempPath)
	if os.IsNotExist(err) {
		//if cannot create it
		if err = os.MkdirAll(config.ServerConf.TempPath, os.ModePerm); err != nil {
			panic(err)
		}
	} else if !stat.IsDir() {
		panic(errors.New("logfile dir is not a directory"))
	}
	logFilePath = filepath.Join(config.ServerConf.TempPath, LogFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	writer = io.MultiWriter(os.Stdout, logFile)
	Logger = hclog.New(&hclog.LoggerOptions{
		Name:   "cfuzz-bot",
		Output: writer,
		Level:  hclog.Debug,
	})

}
