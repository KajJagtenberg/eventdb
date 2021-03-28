package api

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/raft"
)

type RaftService struct {
	raft *raft.Raft
}

func (service *RaftService) Join(ctx context.Context, req *JoinRequest) (*JoinResponse, error) {
	log.Println(req)

	return nil, errors.New("Not implemented")
}

func NewRaftService(raft *raft.Raft) *RaftService {
	return &RaftService{raft}
}
