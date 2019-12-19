package fuzzer

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common/proto"
	"golang.org/x/net/context"
)

type FuzzerGRPCClient struct {
	client proto.FuzzerClient
}

func (f *FuzzerGRPCClient) Prepare(args PrepareArg) error {
	resp, err := f.client.Prepare(context.Background(), &proto.PrepareArg{
		CorpusDir:   args.CorpusDir,
		TargetPath:  args.TargetPath,
		Arguments:   args.Arguments,
		Enviroments: args.Enviroments,
	})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) Fuzz(args FuzzArg) error {
	resp, err := f.client.Fuzz(context.Background(), &proto.FuzzArg{
		TargetPath: args.TargetPath,
		MaxTime:    int32(args.MaxTime),
	})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) Reproduce(args ReproduceArg) error {
	resp, err := f.client.Reproduce(context.Background(), &proto.ReproduceArg{
		TargetPath: args.TargetPath,
		InputPath:  args.InputPath,
		MaxTime:    int32(args.MaxTime),
	})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) MinimizeCorpus(args MinimizeCorpusArg) error {
	resp, err := f.client.MinimizeCorpus(context.Background(), &proto.MinimizeCorpusArg{
		TargetPath: args.TargetPath,
		InputDir:   args.InputDir,
		OutputDir:  args.OutputDir,
		MaxTime:    int32(args.MaxTime),
	})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) Clean() error {
	resp, err := f.client.Clean(context.Background(), &proto.Empty{})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)

}

type FuzzerGRPCServer struct {
	Impl Fuzzer
}

func (f *FuzzerGRPCServer) Prepare(ctx context.Context, args *proto.PrepareArg) (*proto.ErrorResponse, error) {
	err := f.Impl.Prepare(PrepareArg{
		CorpusDir:   args.CorpusDir,
		TargetPath:  args.TargetPath,
		Arguments:   args.Arguments,
		Enviroments: args.Enviroments,
	})
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}

func (f *FuzzerGRPCServer) Fuzz(ctx context.Context, args *proto.FuzzArg) (*proto.ErrorResponse, error) {
	err := f.Impl.Fuzz(FuzzArg{
		TargetPath: args.TargetPath,
		MaxTime:    int(args.MaxTime),
	})
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}

func (f *FuzzerGRPCServer) Reproduce(ctx context.Context, args *proto.ReproduceArg) (*proto.ErrorResponse, error) {
	err := f.Impl.Reproduce(ReproduceArg{
		TargetPath: args.TargetPath,
		InputPath:  args.InputPath,
		MaxTime:    int(args.MaxTime),
	})
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}

func (f *FuzzerGRPCServer) MinimizeCorpus(ctx context.Context, args *proto.MinimizeCorpusArg) (*proto.ErrorResponse, error) {
	err := f.Impl.MinimizeCorpus(MinimizeCorpusArg{
		TargetPath: args.TargetPath,
		InputDir:   args.InputDir,
		OutputDir:  args.OutputDir,
		MaxTime:    int(args.MaxTime),
	})
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}

func (f *FuzzerGRPCServer) Clean(ctx context.Context, args *proto.Empty) (*proto.ErrorResponse, error) {
	err := f.Impl.Clean()
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}
