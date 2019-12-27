package main

import (
	"bytes"
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	//"github.com/go-cmd/cmd"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type AFL struct {
	logger       hclog.Logger
	corpusDir    string
	targetPath   string
	outputDir    string
	arguments    map[string]string
	environments []string
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

func checkAFLOutput(out []string) error {
	for _, v := range out {
		for _, r := range AFL_CHECK_REGEX {
			if r.MatchString(v) {
				return errors.New(COLOR_REGEX.ReplaceAllString(r.FindString(v), ""))
			}
		}
	}
	return nil
}

func (a *AFL) getProgramArg() []string {
	if v, err := a.getArgument(PROGRAM_ARG); err == nil {
		programArgs := strings.Fields(v)
		return programArgs
	}
	return nil
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
	a.environments = args.Environments
	a.corpusDir = args.CorpusDir
	a.targetPath = args.TargetPath

	//check corpus dir and target path exists
	if _, err := os.Stat(args.CorpusDir); os.IsNotExist(err) {
		return errors.New("AFL Prepare CorpusDir not exist")
	}

	if _, err := os.Stat(args.TargetPath); os.IsNotExist(err) {
		return errors.New("AFL Prepare TargetPath not exist")
	}

	outdir, err := ioutil.TempDir(TEMP_DIR, "afl_fuzz")
	if err != nil {
		return errors.New("AFL Prepare create temp directory error: " + err.Error())
	}
	a.outputDir = outdir
	return nil
}

func (a *AFL) Fuzz(args fuzzer.FuzzArg) (fuzzer.FuzzResult, error) {
	a.logger.Debug("fuzz in afl")

	arguments := []string{}
	arguments = append(arguments, INPUT_DIR_FLAG)
	arguments = append(arguments, a.corpusDir)
	arguments = append(arguments, OUTPUT_DIR_FLAG)

	arguments = append(arguments, a.outputDir)

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

	arguments = append(arguments, a.targetPath)
	reproduceArg := []string{a.targetPath}

	//if PROGRAM_ARG exists in arguments, append it to arguments
	programArgs := a.getProgramArg()
	if programArgs != nil {
		arguments = append(arguments, programArgs...)
		reproduceArg = append(reproduceArg, reproduceArg...)
	}

	//run afl
	runner := utils.NewCmd(AFL_PATH, arguments...)
	if len(a.environments) != 0 {
		runner.Env = a.environments
	}

	statusChan := runner.Start()
	//cancelChan := make(chan struct{})

	//a.logger.Debug(AFL_PATH)
	//for _, v := range arguments {
	//	a.logger.Debug(v)
	//}
	//go func() {
	//	ticker := time.NewTicker(time.Duration(AFL_CHECK_TICK_TIME) * time.Second)
	//	for range ticker.C {
	//		status := runner.Status()
	//		for _, v := range status.Stdout {
	//			a.logger.Debug(v)
	//		}
	//	}
	//}()
	go func() {
		<-time.After(time.Duration(args.MaxTime) * time.Second)
		runner.Stop()
	}()

	//finish fuzz
	status := <-statusChan
	err := checkAFLOutput(status.Stdout)
	if err != nil {
		return fuzzer.FuzzResult{}, errors.New("AFL Fuzz error: " + err.Error())
	}

	stats, err := ParseFuzzerStatsFile(a.outputDir)
	if err != nil {
		return fuzzer.FuzzResult{}, errors.New("AFL Fuzz parse stats file error: " + err.Error())
	}

	crashes, err := GetAllCrashes(a.outputDir)
	if err != nil {
		return fuzzer.FuzzResult{}, errors.New("AFL Fuzz get crashes error: " + err.Error())
	}
	for i, _ := range crashes {
		crashes[i].ReproduceArg = reproduceArg
		crashes[i].Environments = a.environments
	}

	return fuzzer.FuzzResult{
		Command:      arguments,
		Crashes:      crashes,
		Stats:        stats,
		TimeExecuted: int(status.Runtime),
	}, nil
}

func (a *AFL) Reproduce(args fuzzer.ReproduceArg) (fuzzer.ReproduceResult, error) {
	a.logger.Debug("reproduce in afl")

	if _, err := os.Stat(args.InputPath); os.IsNotExist(err) {
		return fuzzer.ReproduceResult{}, errors.New("AFL Reproduce InputPath not exist")
	}

	arguments := []string{}
	programArgs := a.getProgramArg()
	isStdinInput := true
	if programArgs != nil {
		for i, v := range programArgs {
			if v == "@@" {
				programArgs[i] = args.InputPath
				isStdinInput = false
			}
		}
		arguments = append(arguments, programArgs...)
	}
	runner := utils.NewCmd(a.targetPath, arguments...)
	if len(a.environments) != 0 {
		runner.Env = a.environments
	}

	if isStdinInput {
		file, err := os.Open(args.InputPath)
		if err != nil {
			return fuzzer.ReproduceResult{}, errors.New("AFL Reproduce cannot open input file: " + err.Error())
		}
		runner.Stdin = file
		defer file.Close()
	}

	statusChan := runner.Start()
	go func() {
		<-time.After(time.Duration(args.MaxTime) * time.Second)
		runner.Stop()
	}()

	//finish reproduce
	status := <-statusChan
	cmd := []string{a.targetPath}
	cmd = append(cmd, programArgs...)

	return fuzzer.ReproduceResult{
		Command:      cmd,
		ReturnCode:   status.Exit,
		TimeExecuted: int(status.Runtime),
		Output:       status.Stdout,
	}, nil
}

func (a *AFL) MinimizeCorpus(args fuzzer.MinimizeCorpusArg) (fuzzer.MinimizeCorpusResult, error) {
	a.logger.Debug("minimize corpus in afl")
	return fuzzer.MinimizeCorpusResult{}, errors.New("AFL MinimizeCorpus Not Implemented")
}

func (a *AFL) Clean() error {
	a.logger.Debug("clean in afl")
	if _, err := os.Stat(a.outputDir); !os.IsNotExist(err) {
		err = os.RemoveAll(a.outputDir)
		if err != nil {
			return errors.New("AFL Clean error :" + err.Error())
		}
	}

	return nil
}
