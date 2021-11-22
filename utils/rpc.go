package utils

import (
	"github.com/CyDrive/consts"
	"github.com/CyDrive/rpc"
)

func PackTransferFileNotification(taskId int32, addr string, filePath string, taskType consts.DataTaskType) *rpc.Notify {
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

func PackDeleteFileNotification(filePath string) *rpc.Notify {
	return &rpc.Notify{
		Notify: &rpc.Notify_DeleteFileNotification{
			DeleteFileNotification: &rpc.DeleteFileNotification{
				FilePath: filePath,
			},
		},
	}
}
