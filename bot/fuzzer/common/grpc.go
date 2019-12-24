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
		CorpusDir:    args.CorpusDir,
		TargetPath:   args.TargetPath,
		Arguments:    args.Arguments,
		Environments: args.Environments,
	})
	if err != nil {
		return err
	}
	if resp.Error == "" {
		return nil
	}
	return errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) Fuzz(args FuzzArg) (FuzzResult, error) {
	resp, err := f.client.Fuzz(context.Background(), &proto.FuzzArg{
		MaxTime: int32(args.MaxTime),
	})
	if err != nil {
		return FuzzResult{}, err
	}

	crashes := []Crash{}
	for _, v := range resp.Crashes {
		crashes = append(crashes, Crash{
			InputPath:    v.InputPath,
			ReproduceArg: v.ReproduceArg,
			Environments: v.Environments,
		})
	}
	result := FuzzResult{
		Command:      resp.Command,
		Crashes:      crashes,
		Stats:        resp.Stats,
		TimeExecuted: int(resp.TimeExecuted),
	}
	if resp.Error == "" {
		return result, nil
	}
	return result, errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) Reproduce(args ReproduceArg) (ReproduceResult, error) {
	resp, err := f.client.Reproduce(context.Background(), &proto.ReproduceArg{
		InputPath: args.InputPath,
		MaxTime:   int32(args.MaxTime),
	})
	if err != nil {
		return ReproduceResult{}, err
	}
	result := ReproduceResult{
		Command:      resp.Command,
		ReturnCode:   int(resp.ReturnCode),
		TimeExecuted: int(resp.TimeExecuted),
		Output:       resp.Output,
	}
	if resp.Error == "" {
		return result, nil
	}
	return result, errors.New(resp.Error)
}

func (f *FuzzerGRPCClient) MinimizeCorpus(args MinimizeCorpusArg) (MinimizeCorpusResult, error) {
	resp, err := f.client.MinimizeCorpus(context.Background(), &proto.MinimizeCorpusArg{
		InputDir:  args.InputDir,
		OutputDir: args.OutputDir,
		MaxTime:   int32(args.MaxTime),
	})
	if err != nil {
		return MinimizeCorpusResult{}, err
	}
	result := MinimizeCorpusResult{
		Command:      resp.Command,
		Stats:        resp.Stats,
		TimeExecuted: int(resp.TimeExecuted),
	}
	if resp.Error == "" {
		return result, nil
	}
	return result, errors.New(resp.Error)
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
		CorpusDir:    args.CorpusDir,
		TargetPath:   args.TargetPath,
		Arguments:    args.Arguments,
		Environments: args.Environments,
	})
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}

func (f *FuzzerGRPCServer) Fuzz(ctx context.Context, args *proto.FuzzArg) (*proto.FuzzResult, error) {
	resp, err := f.Impl.Fuzz(FuzzArg{
		MaxTime: int(args.MaxTime),
	})
	crashes := []*proto.Crash{}
	for _, v := range resp.Crashes {
		crashes = append(crashes, &proto.Crash{
			InputPath:    v.InputPath,
			ReproduceArg: v.ReproduceArg,
			Environments: v.Environments,
		})
	}
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return &proto.FuzzResult{
		Command:      resp.Command,
		Crashes:      crashes,
		Stats:        resp.Stats,
		TimeExecuted: int32(resp.TimeExecuted),
		Error:        errMsg,
	}, nil
}

func (f *FuzzerGRPCServer) Reproduce(ctx context.Context, args *proto.ReproduceArg) (*proto.ReproduceResult, error) {
	resp, err := f.Impl.Reproduce(ReproduceArg{
		InputPath: args.InputPath,
		MaxTime:   int(args.MaxTime),
	})
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return &proto.ReproduceResult{
		Command:      resp.Command,
		ReturnCode:   int32(resp.ReturnCode),
		TimeExecuted: int32(resp.TimeExecuted),
		Output:       resp.Output,
		Error:        errMsg,
	}, nil
}

func (f *FuzzerGRPCServer) MinimizeCorpus(ctx context.Context, args *proto.MinimizeCorpusArg) (*proto.MinimizeCorpusResult, error) {
	resp, err := f.Impl.MinimizeCorpus(MinimizeCorpusArg{
		InputDir:  args.InputDir,
		OutputDir: args.OutputDir,
		MaxTime:   int(args.MaxTime),
	})
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return &proto.MinimizeCorpusResult{
		Command:      resp.Command,
		Stats:        resp.Stats,
		TimeExecuted: int32(resp.TimeExecuted),
		Error:        errMsg,
	}, nil
}

func (f *FuzzerGRPCServer) Clean(ctx context.Context, args *proto.Empty) (*proto.ErrorResponse, error) {
	err := f.Impl.Clean()
	if err == nil {
		return &proto.ErrorResponse{Error: ""}, nil
	}
	return &proto.ErrorResponse{Error: err.Error()}, nil
}
