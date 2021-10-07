package master

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "../rpc"
	"github.com/yah01/CyDrive/master/node_manager"
)

func (m *Master) JoinCluster(ctx context.Context, req *pb.JoinClusterRequest) (resp *pb.JoinClusterResponse, error) {
	nodeManager := node_manager.GetNodeManager()
	
	node := node_manager.NewNode(req.Capacity, req.Usage)

	nodeManager.AddNode(node)

	resp = &pb.JoinClusterResponse{
		Id: node.Id,
	}

	return resp, nil
}

func (m *Master) HeartBeats(ctx context.Context, req *pb.HeartBeatsRequest) (resp *pb.HeartBeatsResponse, error) {
	nodeManager := node_manager.GetNodeManager()

	node := nodeManager.GetNode(req.Id)
	if node == nil {
		return nil,fmt.Errorf("no such node, join cluster first")
	}

	node.LastHeartBeatTime = time.Now()

	resp = &pb.HeartBeatsResponse{}
	return resp,nil
}
