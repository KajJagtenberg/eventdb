package cluster

import "context"

type ClusterService struct {
	cl *Cluster
}

func (s *ClusterService) GetHealthScore(ctx context.Context, req *GetHealthScoreRequest) (*GetHealthScoreResponse, error) {
	return nil, nil // TODO: Implement
}

func (s *ClusterService) GetNodes(ctx context.Context, req *GetNodesRequest) (*GetNodesResponse, error) {
	return nil, nil // TODO: Implement
}

func NewClusterService(cl *Cluster) *ClusterService {
	return &ClusterService{cl}
}
