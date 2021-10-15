package master

import (
	"context"
	"fmt"
	"time"

	"github.com/CyDrive/master/node_manager"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/utils"
)

type NodeManageServer struct {
	rpc.UnimplementedManageServer
}

func (s *NodeManageServer) JoinCluster(ctx context.Context, req *rpc.JoinClusterRequest) (*rpc.JoinClusterResponse, error) {
	nodeManager := GetNodeManager()
	node := node_manager.NewNode(req.Capacity, req.Usage)

	nodeManager.AddNode(node)

	resp := &rpc.JoinClusterResponse{
		Id: node.Id,
	}

	return resp, nil
}

func (s *NodeManageServer) HeartBeats(ctx context.Context, req *rpc.HeartBeatsRequest) (*rpc.HeartBeatsResponse, error) {
	nodeManager := GetNodeManager()

	node := nodeManager.GetNode(req.Id)
	if node == nil {
		return nil, fmt.Errorf("no such node, join cluster first")
	}

	node.LastHeartBeatTime = time.Now()

	resp := &rpc.HeartBeatsResponse{}
	return resp, nil
}

func (s *NodeManageServer) Notifier(req *rpc.ConnectNotifierRequest, stream rpc.Manage_NotifierServer) error {
	nodeManager := GetNodeManager()
	notifyChan, ok := nodeManager.GetNotifyChan(req.NodeId)
	if !ok {
		return fmt.Errorf("no such node, nodeId=%d, node may haven't join the cluster", req.NodeId)
	}

	for notificationI := range notifyChan {
		switch notification := notificationI.(type) {
		case *rpc.CreateSendFileTaskNotify:
			stream.Send(utils.PackCreateSendFileTaskNotify(notification))

		case *rpc.CreateRecvFileTaskNotify:
			stream.Send(utils.PackCreateRecvFileTaskNotify(notification))
		}
	}

	return nil
}
