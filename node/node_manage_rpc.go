package node

import (
	"context"
	"time"

	"github.com/CyDrive/rpc"
	log "github.com/sirupsen/logrus"
)

func (node *StorageNode) JoinCluster() error {
	req := &rpc.JoinClusterRequest{
		Capacity: node.Cap,
		Usage:    node.Usage,
		Type:     node.Type,
	}

	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	resp, err := node.manageClient.JoinCluster(ctx, req)
	if err != nil {
		return err
	}

	node.Id = resp.Id

	log.Infof("join cluster, assigned id: %d", node.Id)
	return nil
}

func (node *StorageNode) HeartBeat() {
	req := &rpc.HeartBeatsRequest{
		Id:              node.Id,
		StorageUsage:    node.Usage,
		CpuUsagePercent: 0,
		TaskNum:         0,
	}

	ctx, _ := context.WithTimeout(context.Background(), 250*time.Millisecond)
	_, err := node.manageClient.HeartBeats(ctx, req)
	if err != nil {
		log.Errorf("heartbeat failed, err=%+v", err)
		return
	}

	log.Infof("heartbeat with req=%+v", req)
}

func (node *StorageNode) Notify() {
	stream, err := node.manageClient.Notifier(context.Background(), &rpc.ConnectNotifierRequest{
		NodeId: node.Id,
	})
	if err != nil {

	}

	notify, err := stream.Recv()

	switch notification := notify.Notify.(type) {
	case *rpc.Notify_CreateSendfileTaskNotify:
		sendFileNotify := notification.CreateSendfileTaskNotify
		log.Infof("notify=%+v", sendFileNotify)
	}
}
