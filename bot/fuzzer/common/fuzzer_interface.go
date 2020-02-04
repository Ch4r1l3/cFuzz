package fuzzer

import (
	"context"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"net/rpc"
)

type PrepareArg struct {
	CorpusDir    string //directory of corpus
	TargetPath   string
	Arguments    map[string]string
	Environments []string
}

type FuzzArg struct {
	MaxTime int
}

type Crash struct {
	InputPath    string //Path of the crash file
	FileName     string // The origin filename, this might different from the filename in inputpath
	ReproduceArg []string
	Environments []string
}

type FuzzResult struct {
	Command      []string
	Crashes      []Crash
	Stats        map[string]string
	TimeExecuted int
}

type ReproduceArg struct {
	InputPath string //Path of the input file to be reproduced
	MaxTime   int
}

type ReproduceResult struct {
	Command      []string
	ReturnCode   int
	TimeExecuted int
	Output       []string
}

type Fuzzer interface {
	Prepare(args PrepareArg) error
	Fuzz(args FuzzArg) (FuzzResult, error)
	Reproduce(args ReproduceArg) (ReproduceResult, error)
	Clean() error
}

type FuzzerPlugin struct {
	Impl Fuzzer
}

func (f *FuzzerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FuzzerRPCServer{Impl: f.Impl}, nil
}

func (FuzzerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FuzzerRPCClient{client: c}, nil
}

type FuzzerGRPCPlugin struct {
	plugin.Plugin
	Impl Fuzzer
}

func (f *FuzzerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterFuzzerServer(s, &FuzzerGRPCServer{Impl: f.Impl})
	return nil
}

func (f *FuzzerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &FuzzerGRPCClient{client: proto.NewFuzzerClient(c)}, nil
}
