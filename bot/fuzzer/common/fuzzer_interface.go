package fuzzer

import (
	"./proto"
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"net/rpc"
)

type Fuzzer interface {
	Prepare(corpus_dir string, target_path string, arguments map[string]string, enviroments map[string]string) error
	Fuzz(target_path string, max_time int) error
	Reproduce(target_path string, input_path string, max_time int) error
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
