package fuzzer

import (
	"./proto"
	"errors"
	"golang.org/x/net/context"
)

type FuzzerGRPCClient struct {
	client proto.FuzzerClient
}

func (f *FuzzerGRPCClient) Fuzz() error {
	resp, err := f.client.Fuzz(context.Background(), &proto.Empty{})
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

type FuzzerGRPCServer struct {
	Impl Fuzzer
}

func (f *FuzzerGRPCServer) Fuzz(ctx context.Context, emp *proto.Empty) (*proto.FuzzResponse, error) {
	err := f.Impl.Fuzz()
	if err == nil {
		return &proto.FuzzResponse{Error: ""}, nil
	}
	return &proto.FuzzResponse{Error: err.Error()}, nil
}
