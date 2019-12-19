package main

import (
	"bytes"
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/go-cmd/cmd"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type AFL struct {
	logger      hclog.Logger
	corpusDir   string
	targetPath  string
	arguments   map[string]string
	enviroments []string
}

func (a *AFL) getArgument(key string) (string, error) {
	if a.arguments == nil {
		return "", errors.New("arguments is nil")
	}
	if v, found := a.arguments[key]; found {
		return v, nil
	}
	return "", errors.New("key not in arguments")
}

func (a *AFL) Prepare(args fuzzer.PrepareArg) error {
	a.logger.Debug("prepare in afl")

	//check core pattern
	f, err := os.Open(CORE_PATTERN_FILE_PATH)
	if err != nil {
		return errors.New("AFL Prepare open " + CORE_PATTERN_FILE_PATH + " error: " + err.Error())
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	content := strings.TrimSpace(buf.String())

	if content != "core" {
		//core pattern file content not equal to core, fix it
		bcore := []byte("core")
		err = ioutil.WriteFile(CORE_PATTERN_FILE_PATH, bcore, 0644)
		if err != nil {
			return errors.New("AFL Prepare fix core pattern file fail error: " + err.Error())
		}
	}

	//save args
	a.arguments = args.Arguments
	a.enviroments = args.Enviroments
	a.corpusDir = args.CorpusDir
	a.targetPath = args.TargetPath

	//check corpus dir and target path exists
	if _, err := os.Stat(args.CorpusDir); os.IsNotExist(err) {
		return errors.New("AFL Prepare CorpusDir not exist")
	}

	if _, err := os.Stat(args.TargetPath); os.IsNotExist(err) {
		return errors.New("AFL Prepare TargetPath not exist")
	}

	return nil
}

func (a *AFL) Fuzz(args fuzzer.FuzzArg) error {
	a.logger.Debug("fuzz in afl")
	if a.targetPath != "" && a.targetPath != args.TargetPath {
		return errors.New("AFL Fuzz TargetPath not equal to the TargetPath in Prepare")
	}

	arguments := []string{}
	arguments = append(arguments, INPUT_DIR_FLAG)
	arguments = append(arguments, a.corpusDir)
	arguments = append(arguments, OUTPUT_DIR_FLAG)

	dir, err := ioutil.TempDir(TEMP_DIR, "afl_fuzz")
	if err != nil {
		return errors.New("AFL Fuzz create temp directory error: " + err.Error())
	}
	arguments = append(arguments, dir)

	//if MEMORY_LIMIT exists in arguments, append it to arguments
	if v, err := a.getArgument(MEMORY_LIMIT); err == nil {
		arguments = append(arguments, MEMORY_LIMIT_FLAG)
		arguments = append(arguments, v)
	}

	//if TIMEOUT_LIMIT exists in arguments, append it to arguments
	if v, err := a.getArgument(TIMEOUT_LIMIT); err == nil {
		arguments = append(arguments, TIMEOUT_LIMIT_FLAG)
		arguments = append(arguments, v)
	}

	arguments = append(arguments, args.TargetPath)

	//if PROGRAM_ARG exists in arguments, append it to arguments
	if v, err := a.getArgument(PROGRAM_ARG); err == nil {
		arguments = append(arguments, v)
	}

	//run afl
	runner := cmd.NewCmd(AFL_PATH, arguments...)
	if len(a.enviroments) != 0 {
		runner.Env = a.enviroments
	}

	statusChan := runner.Start()

	go func() {
		<-time.After(time.Duration(args.MaxTime) * time.Second)
		runner.Stop()
	}()

	//finish fuzz
	<-statusChan

	return nil
}

func (a *AFL) Reproduce(args fuzzer.ReproduceArg) error {
	a.logger.Debug("reproduce in afl")
	return nil
}

func (a *AFL) MinimizeCorpus(args fuzzer.MinimizeCorpusArg) error {
	a.logger.Debug("minimize corpus in afl")
	return nil
}

func (a *AFL) Clean() error {
	a.logger.Debug("clean in afl")
	return nil
}
