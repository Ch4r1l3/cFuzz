package main

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	//"github.com/go-cmd/cmd"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type LibFuzzer struct {
	logger       hclog.Logger
	corpusDir    string
	newCorpusDir string
	targetPath   string
	outputDir    string
	arguments    map[string]string
	environments []string
	mergeSuccess bool
}

func (a *LibFuzzer) getArgument(key string) (string, error) {
	if a.arguments == nil {
		return "", errors.New("arguments is nil")
	}
	if v, found := a.arguments[key]; found {
		return v, nil
	}
	return "", errors.New("key not in arguments")
}

func (a *LibFuzzer) getProgramArg() []string {
	if v, err := a.getArgument(PROGRAM_ARG); err == nil {
		programArgs := strings.Fields(v)
		return programArgs
	}
	return nil
}

func (a *LibFuzzer) Prepare(args fuzzer.PrepareArg) error {
	a.logger.Debug("prepare in libFuzzer")

	//save args
	a.arguments = args.Arguments
	a.environments = args.Environments
	a.corpusDir = args.CorpusDir
	a.targetPath = args.TargetPath

	//check corpus dir and target path exists
	if _, err := os.Stat(args.CorpusDir); os.IsNotExist(err) {
		return errors.New("LibFuzzer Prepare CorpusDir not exist")
	}

	if _, err := os.Stat(args.TargetPath); os.IsNotExist(err) {
		return errors.New("LibFuzzer Prepare TargetPath not exist")
	}

	outdir, err := ioutil.TempDir(TEMP_DIR, "libFuzzer_fuzz")
	if err != nil {
		return errors.New("LibFuzzer Prepare create temp directory error: " + err.Error())
	}
	a.outputDir = outdir

	a.newCorpusDir, err = ioutil.TempDir(TEMP_DIR, "libFuzzer_corpus")
	if err != nil {
		return errors.New("LibFuzzer Prepare create temp directory error: " + err.Error())
	}
	runner := utils.NewCmd(a.targetPath, MERGE, a.newCorpusDir, a.corpusDir)
	a.mergeSuccess = true
	statusChan := runner.Start()
	finishChan := make(chan struct{})

	go func() {
		select {
		case <-time.After(time.Duration(MAX_MERGE_TIME) * time.Second):
			runner.Stop()
			a.mergeSuccess = false
		case <-finishChan:
			return
		}
	}()

	<-statusChan
	finishChan <- struct{}{}

	return nil
}

func (a *LibFuzzer) Fuzz(args fuzzer.FuzzArg) (fuzzer.FuzzResult, error) {
	a.logger.Debug("fuzz in libFuzzer")

	arguments := []string{}

	// get final stats
	arguments = append(arguments, FINAL_STATS)

	reproduceArg := []string{a.targetPath}

	// if PROGRAM_ARG exists in arguments, append it to arguments
	programArgs := a.getProgramArg()
	if programArgs != nil {
		arguments = append(arguments, programArgs...)
		reproduceArg = append(reproduceArg, programArgs...)
	}

	// append corpusDir
	arguments = append(arguments, a.newCorpusDir)
	if !a.mergeSuccess {
		arguments = append(arguments, a.corpusDir)
	}

	// max total time
	arguments = append(arguments, MAX_TOTAL_TIME+strconv.Itoa(args.MaxTime-1))

	// ouput dir
	outputDir := OUTPUT_FLAG + a.outputDir
	if outputDir[len(outputDir)-1:] != "/" {
		outputDir += "/"
	}
	arguments = append(arguments, outputDir)

	// run libFuzzer
	runner := utils.NewCmd(a.targetPath, arguments...)
	if len(a.environments) != 0 {
		runner.Env = a.environments
	}

	statusChan := runner.Start()
	// cancelChan := make(chan struct{})

	// a.logger.Debug(LibFuzzer_PATH)
	// for _, v := range arguments {
	// 	a.logger.Debug(v)
	// }
	// go func() {
	// 	ticker := time.NewTicker(time.Duration(LibFuzzer_CHECK_TICK_TIME) * time.Second)
	// 	for range ticker.C {
	// 		status := runner.Status()
	// 		for _, v := range status.Stdout {
	// 			a.logger.Debug(v)
	// 		}
	// 	}
	// }()
	go func() {
		<-time.After(time.Duration(args.MaxTime) * time.Second)
		runner.Stop()
	}()

	//finish fuzz
	status := <-statusChan

	stats, err := ParseFuzzerStats(status.Stdout)
	if err != nil {
		return fuzzer.FuzzResult{}, errors.New("LibFuzzer Fuzz parse stats error: " + err.Error())
	}

	crashes, err := a.GetAllCrashes()
	if err != nil {
		return fuzzer.FuzzResult{}, errors.New("LibFuzzer Fuzz get crashes error: " + err.Error())
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

func (a *LibFuzzer) Reproduce(args fuzzer.ReproduceArg) (fuzzer.ReproduceResult, error) {
	a.logger.Debug("reproduce in libFuzzer")

	if _, err := os.Stat(args.InputPath); os.IsNotExist(err) {
		return fuzzer.ReproduceResult{}, errors.New("LibFuzzer Reproduce InputPath not exist")
	}

	arguments := []string{}
	programArgs := a.getProgramArg()
	arguments = append(arguments, programArgs...)
	arguments = append(arguments, args.InputPath)

	runner := utils.NewCmd(a.targetPath, arguments...)
	if len(a.environments) != 0 {
		runner.Env = a.environments
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

func (a *LibFuzzer) Clean() error {
	a.logger.Debug("clean in libFuzzer")
	if _, err := os.Stat(a.outputDir); !os.IsNotExist(err) {
		err = os.RemoveAll(a.outputDir)
		if err != nil {
			return errors.New("LibFuzzer Clean error :" + err.Error())
		}
	}

	return nil
}
