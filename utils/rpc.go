package utils

import (
	"github.com/CyDrive/consts"
	"github.com/CyDrive/rpc"
)

func PackCreateTransferFileTaskNotification(taskId int32, addr string, filePath string, taskType consts.DataTaskType) *rpc.Notify {
	return &rpc.Notify{
		Notify: &rpc.Notify_TransferFileNotification{
			TransferFileNotification: &rpc.TransferFileNotification{
				TaskId:   taskId,
				Addr:     addr,
				FilePath: filePath,
				TaskType: taskType,
			},
		},
	}
}
