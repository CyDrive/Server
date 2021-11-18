package utils

import "github.com/CyDrive/rpc"

func PackCreateTransferFileTaskNotification(taskId int32, addr string, filePath string) *rpc.Notify {
	return &rpc.Notify{
		Notify: &rpc.Notify_CreateFileTransferTaskNotification{
			CreateFileTransferTaskNotification: &rpc.CreateFileTransferTaskNotification{
				TaskId:   taskId,
				Addr:     addr,
				FilePath: filePath,
			},
		},
	}
}
