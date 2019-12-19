package main

import (
	"bytes"
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"os"
	"strings"
)

type AFL struct {
	logger      hclog.Logger
	corpusDir   string
	targetPath  string
	arguments   map[string]string
	enviroments map[string]string
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
	if a.targetPath != "" && a.targetPath != args.TargetPath {
		return errors.New("AFL Fuzz TargetPath not equal to the TargetPath in Prepare")
	}
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
