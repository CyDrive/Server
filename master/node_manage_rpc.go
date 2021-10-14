package master

import (
	"context"
	"fmt"
	"time"

	"github.com/CyDrive/master/node_manager"
	"github.com/CyDrive/rpc"
)

func (s *NodeManagerServer) JoinCluster(ctx context.Context, req *rpc.JoinClusterRequest) (*rpc.JoinClusterResponse, error) {
	nodeManager := node_manager.GetNodeManager()
	node := node_manager.NewNode(req.Capacity, req.Usage)

	nodeManager.AddNode(node)

	resp := &rpc.JoinClusterResponse{
		Id: node.Id,
	}

	return resp, nil
}

func (s *NodeManagerServer) HeartBeats(ctx context.Context, req *rpc.HeartBeatsRequest) (*rpc.HeartBeatsResponse, error) {
	nodeManager := node_manager.GetNodeManager()

	node := nodeManager.GetNode(req.Id)
	if node == nil {
		return nil, fmt.Errorf("no such node, join cluster first")
	}

	node.LastHeartBeatTime = time.Now()

	resp := &rpc.HeartBeatsResponse{}
	return resp, nil
}
