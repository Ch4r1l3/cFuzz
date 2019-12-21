package fuzzer

import (
	"context"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"net/rpc"
)

type PrepareArg struct {
	CorpusDir   string
	TargetPath  string
	Arguments   map[string]string
	Enviroments []string
}

type FuzzArg struct {
	MaxTime int
}

type Crash struct {
	InputPath    string //Path of the crash file
	ReproduceArg []string
	Enviroments  []string
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

type MinimizeCorpusArg struct {
	InputDir  string //Path of the corpus to be minimized
	OutputDir string
	MaxTime   int
}

type MinimizeCorpusResult struct {
	Command      []string
	Stats        map[string]string
	TimeExecuted int
}

type Fuzzer interface {
	Prepare(args PrepareArg) error
	Fuzz(args FuzzArg) (FuzzResult, error)
	Reproduce(args ReproduceArg) (ReproduceResult, error)
	MinimizeCorpus(args MinimizeCorpusArg) (MinimizeCorpusResult, error)
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
