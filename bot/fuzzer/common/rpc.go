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
	f.client.Call("Plugin.Fuzz", args, &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

func (f *FuzzerRPCClient) Fuzz(args FuzzArg) error {
	var resp string
	f.client.Call("Plugin.Fuzz", args, &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

func (f *FuzzerRPCClient) Reproduce(args ReproduceArg) error {
	var resp string
	f.client.Call("Plugin.Reproduce", args, &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

func (f *FuzzerRPCClient) MinimizeCorpus(args MinimizeCorpusArg) error {
	var resp string
	f.client.Call("Plugin.MinimizeCorpus", args, &resp)
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

func (f *FuzzerRPCServer) Fuzz(args FuzzArg, resp *string) error {
	err := f.Impl.Fuzz(args)
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return nil
}

func (f *FuzzerRPCServer) Reproduce(args ReproduceArg, resp *string) error {
	err := f.Impl.Reproduce(args)
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return err
}

func (f *FuzzerRPCServer) MinimizeCorpus(args MinimizeCorpusArg, resp *string) error {
	err := f.Impl.MinimizeCorpus(args)
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return err
}
