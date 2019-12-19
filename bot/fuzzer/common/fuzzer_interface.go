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
	TargetPath string
	MaxTime    int
}

type ReproduceArg struct {
	TargetPath string
	InputPath  string
	MaxTime    int
}

type MinimizeCorpusArg struct {
	TargetPath string
	InputDir   string
	OutputDir  string
	MaxTime    int
}

type Fuzzer interface {
	Prepare(args PrepareArg) error
	Fuzz(args FuzzArg) error
	Reproduce(args ReproduceArg) error
	MinimizeCorpus(args MinimizeCorpusArg) error
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
