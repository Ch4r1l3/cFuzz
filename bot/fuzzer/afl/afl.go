package main

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/hashicorp/go-hclog"
)

type AFL struct {
	logger hclog.Logger
}

func (a *AFL) Prepare(args fuzzer.PrepareArg) error {
	a.logger.Debug("prepare in afl")
	return nil
}

func (a *AFL) Fuzz(args fuzzer.FuzzArg) error {
	a.logger.Debug("fuzz in afl")
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
