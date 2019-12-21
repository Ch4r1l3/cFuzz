package fuzzer

import (
	"errors"
	"net/rpc"
)

type FuzzerRPCClient struct {
	client *rpc.Client
}

func (f *FuzzerRPCClient) Prepare(args PrepareArg) error {
	var resp string
	f.client.Call("Plugin.Prepare", args, &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

func (f *FuzzerRPCClient) Fuzz(args FuzzArg) (FuzzResult, error) {
	var resp FuzzResult
	err := f.client.Call("Plugin.Fuzz", args, &resp)
	return resp, err
}

func (f *FuzzerRPCClient) Reproduce(args ReproduceArg) (ReproduceResult, error) {
	var resp ReproduceResult
	err := f.client.Call("Plugin.Reproduce", args, &resp)
	return resp, err
}

func (f *FuzzerRPCClient) MinimizeCorpus(args MinimizeCorpusArg) (MinimizeCorpusResult, error) {
	var resp MinimizeCorpusResult
	err := f.client.Call("Plugin.MinimizeCorpus", args, &resp)
	return resp, err
}

func (f *FuzzerRPCClient) Clean() error {
	var resp string
	f.client.Call("Plugin.Clean", new(interface{}), &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

type FuzzerRPCServer struct {
	Impl Fuzzer
}

func (f *FuzzerRPCServer) Prepare(args PrepareArg, resp *string) error {
	err := f.Impl.Prepare(args)
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return err
}

func (f *FuzzerRPCServer) Fuzz(args FuzzArg, resp *FuzzResult) error {
	v, err := f.Impl.Fuzz(args)
	*resp = v
	return err
}

func (f *FuzzerRPCServer) Reproduce(args ReproduceArg, resp *ReproduceResult) error {
	v, err := f.Impl.Reproduce(args)
	*resp = v
	return err
}

func (f *FuzzerRPCServer) MinimizeCorpus(args MinimizeCorpusArg, resp *MinimizeCorpusResult) error {
	v, err := f.Impl.MinimizeCorpus(args)
	*resp = v
	return err
}

func (f *FuzzerRPCServer) Clean(args interface{}, resp *string) error {
	err := f.Impl.Clean()
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return err
}
