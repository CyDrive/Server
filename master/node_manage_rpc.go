package master

import (
	"context"
	"fmt"
	"time"

	"github.com/CyDrive/master/managers"
	"github.com/CyDrive/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

type NodeManageServer struct {
	rpc.UnimplementedManageServer
}

func (s *NodeManageServer) JoinCluster(ctx context.Context, req *rpc.JoinClusterRequest) (*rpc.JoinClusterResponse, error) {
	nodeManager := GetNodeManager()
	peer, _ := peer.FromContext(ctx)
	node := managers.NewNode(req.Capacity, req.Usage, peer.Addr.String())

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

	for notification := range notifyChan {
		switch notify := notification.Notify.(type) {
		case *rpc.Notify_CreateFileTransferTaskNotification:
			log.Infof("notify node=%v, notification=%+v", req.NodeId, notify)
			err := stream.Send(&rpc.Notify{
				Notify: notify,
			})
			if err != nil {
				log.Errorf("failed to notify, err=%v", err)
				break
			}
		}
	}

	return nil
}
