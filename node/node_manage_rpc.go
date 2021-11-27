package node

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/utils"
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
		TaskNum:         node.TaskNum,
		State:           node.State,
	}

	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	_, err := node.manageClient.HeartBeats(ctx, req)
	if err != nil {
		log.Errorf("heartbeat failed, err=%+v", err)
		return
	}

	// log.Infof("heartbeat with req=%+v", req)
}

func (node *StorageNode) ProcessNotifications() {
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
			continue
		}

		log.Infof("recv delete file notification: %+v", notify.Notify)

		switch notify := notify.Notify.(type) {
		case *rpc.Notify_TransferFileNotification:
			notification := notify.TransferFileNotification

			node.TaskNum++
			if notification.TaskType == consts.DataTaskType_Download {
				go node.DownloadFile(notification.TaskId, notification.FilePath, notification.Addr+consts.FileTransferorListenPortStr)
			} else if notification.TaskType == consts.DataTaskType_Upload {
				go node.UploadFile(notification.TaskId, notification.FilePath, notification.Addr+consts.FileTransferorListenPortStr)
			}

		case *rpc.Notify_DeleteFileNotification:
			notification := notify.DeleteFileNotification
			err = os.RemoveAll(notification.FilePath)
			if err != nil {
				log.Errorf("failed to remove file for path=%s, err=%+v", notification.FilePath, err)
			}
		}
	}
}

func (node *StorageNode) ReportFileInfos() {
	ctx := context.Background()
	req := &rpc.ReportFileInfosRequest{
		Id:        node.Id,
		FileInfos: make([]*models.FileInfo, 0, 10),
	}
	err := filepath.Walk(node.StorePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.Trim(strings.ReplaceAll(path, "\\", "/"), "/")

		req.FileInfos = append(req.FileInfos, utils.NewFileInfo(info, path))
		if len(req.FileInfos) >= consts.ReportFileInfoBatchSize {
			_, err = node.manageClient.ReportFileInfos(ctx, req)
			if err != nil {
				return err
			}
			log.Infof("report fileInfos: %+v", req.FileInfos)
			req.FileInfos = req.FileInfos[:0]
		}

		return nil
	})

	if err == nil && len(req.FileInfos) > 0 { // some file infos left
		_, err = node.manageClient.ReportFileInfos(ctx, req)
		if err == nil {
			log.Infof("report fileInfos: %+v", req.FileInfos)
			req.FileInfos = req.FileInfos[:0]
		}
	}

	if err != nil {
		log.Errorf("failed to report file infos, err=%v, fileInfos=%+v", err, req.FileInfos)
		return
	}
	node.State = consts.NodeState_Running
}
