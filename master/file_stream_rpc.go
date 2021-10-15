package master

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/CyDrive/master/envs"
	"github.com/CyDrive/network"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/types"
)

type FileStreamServer struct {
	rpc.UnimplementedFileStreamServer
}

func (s *FileStreamServer) SendFile(stream rpc.FileStream_SendFileServer) error {
	var task *network.DataTask = nil

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			task.FileHandle.Close()
			return stream.SendAndClose(&rpc.SendFileResponse{})
		}
		if err != nil {
			return err
		}

		if len(req.Error) > 0 {
			task.FileHandle.(*envs.RemoteFile).Err = errors.New(req.Error)
			return nil
		}

		if task == nil {
			if req.TaskId > 0 {
				task = GetFileTransferor().
					GetTask(types.TaskId(req.TaskId))
			} else {
				task.FileHandle.(*envs.RemoteFile).Err =
					fmt.Errorf("node sends no task_id")

				return fmt.Errorf("need task_id, req=%+v", req)
			}
		}

		io.Copy(task.FileHandle, bytes.NewReader(req.Data))
	}
}
