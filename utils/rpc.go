package utils

import "github.com/CyDrive/rpc"

func PackCreateSendFileTaskNotify(req *rpc.CreateSendFileTaskNotify) *rpc.Notify {
	return &rpc.Notify{
		Notify: &rpc.Notify_CreateSendfileTaskNotify{
			CreateSendfileTaskNotify: req,
		},
	}
}

func PackCreateRecvFileTaskNotify(req *rpc.CreateRecvFileTaskNotify) *rpc.Notify {
	return &rpc.Notify{
		Notify: &rpc.Notify_CreateRecvfileTaskNotify{
			CreateRecvfileTaskNotify: req,
		},
	}
}
