package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ShellService struct{}

func (service *ShellService) Execute(ctx context.Context, req *ShellRequest) (*ShellResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func NewShellService() *ShellService {
	return &ShellService{}
}
