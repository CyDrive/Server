package node

import (
	"context"
	"time"

	"github.com/CyDrive/consts"
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
		log.Errorf("failed to connect the notifier, err=%v", err)
		panic("failed to connect the notifier")
	}

	for {
		notify, err := stream.Recv()
		if err != nil {
			log.Errorf("failed to recv notification, err=%v", err)
		}

		switch notify := notify.Notify.(type) {
		case *rpc.Notify_CreateFileTransferTaskNotification:
			notification := notify.CreateFileTransferTaskNotification
			log.Infof("recv file transfer notification: %+v", notification)
			if notification.TaskType == consts.DataTaskType_Download {
				node.DownloadFile(notification.TaskId, node.StorePath+"/"+notification.FilePath, notification.Addr+consts.FileTransferorListenPortStr)
			} else if notification.TaskType == consts.DataTaskType_Upload {
				node.UploadFile(notification.TaskId, node.StorePath+"/"+notification.FilePath, notification.Addr+consts.FileTransferorListenPortStr)

			}
		}
	}

}
