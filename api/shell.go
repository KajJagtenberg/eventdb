package api

import (
	"io"

	"github.com/kajjagtenberg/eventflowdb/shell"
)

type ShellService struct {
}

func (service *ShellService) Execute(stream ShellService_ExecuteServer) error {
	shell, err := shell.NewShell()
	if err != nil {
		return err
	}

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		result, err := shell.Execute(request.Body)
		if err != nil {
			return err
		}

		if err := stream.Send(&ShellResponse{
			Body: result,
		}); err != nil {
			return err
		}
	}
}

func NewShellService() *ShellService {
	return &ShellService{}
}
