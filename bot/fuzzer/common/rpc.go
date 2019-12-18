package fuzzer

import (
	"errors"
	"net/rpc"
)

type FuzzerRPCClient struct {
	client *rpc.Client
}

func (f *FuzzerRPCClient) Fuzz() error {
	var resp string
	f.client.Call("Plugin.Fuzz", new(interface{}), &resp)
	if resp == "" {
		return nil
	}
	return errors.New(resp)
}

type FuzzerRPCServer struct {
	Impl Fuzzer
}

func (f *FuzzerRPCServer) Fuzz(args interface{}, resp *string) error {
	err := f.Impl.Fuzz()
	if err != nil {
		*resp = err.Error()
	} else {
		*resp = ""
	}
	return nil
}
