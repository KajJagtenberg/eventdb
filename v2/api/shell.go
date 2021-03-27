package api

import (
	"github.com/hashicorp/raft"
	"github.com/kajjagtenberg/eventflowdb/shell"
)

type ShellService struct {
	raft *raft.Raft
}

func (service *ShellService) Execute(stream ShellService_ExecuteServer) error {
	shell := shell.NewShell(service.raft)

	for {
		request, err := stream.Recv()
		if err != nil {
			return err
		}

		output, err := shell.Execute(request.Body)

		var body string
		if err == nil {
			body = output
		} else {
			body = err.Error()
		}

		if err := stream.Send(&ShellResponse{
			Body: body,
		}); err != nil {
			return err
		}
	}
}

func NewShellService(raft *raft.Raft) *ShellService {
	return &ShellService{raft}
}
