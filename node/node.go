package node

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/node/config"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type StorageNode struct {
	*config.Config
	Id      int32
	Usage   int64
	TaskNum int32

	heartBeatTimer   *time.Timer
	conn             *grpc.ClientConn
	manageClient     rpc.ManageClient
	fileStreamClient rpc.FileStreamClient
}

func NewStorageNode(config *config.Config) *StorageNode {
	usage, err := utils.DirSize(config.StorePath)
	if err != nil {
		panic(err)
	}

	logfile, err := os.Create(fmt.Sprintf("node %v.log", utils.GetDateTimeNow()))
	if err != nil {
		panic(err)
	}

	log.SetOutput(logfile)

	node := StorageNode{
		Config: config,
		Usage:  usage,

		heartBeatTimer: time.NewTimer(250 * time.Millisecond),
	}

	return &node
}

func (node *StorageNode) Start() {
	var err error

	log.Infof("connecting to the gRPC services")
	node.conn, err = grpc.Dial(node.MasterAddr+consts.RpcListenPortStr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	node.manageClient = rpc.NewManageClient(node.conn)
	node.fileStreamClient = rpc.NewFileStreamClient(node.conn)

	log.Infof("connections setups...")
	// join cluster
	for {
		if err := node.JoinCluster(); err != nil {
			log.Errorf("can't join cluster, err=%+v, will retry after 500ms", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	// cron tasks
	go func() {
		for {
			select {
			case <-node.heartBeatTimer.C:
				node.heartBeatTimer.Reset(250 * time.Millisecond)
				node.HeartBeat()
			}
		}
	}()

	go node.ReportFileInfos()

	// process notifications
	node.ProcessNotifications()
}

func (node *StorageNode) DownloadFile(taskId types.TaskId, filePath, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("failed to connect to the peer, addr=%v, err=%v", addr, err)
		return
	}
	defer conn.Close()

	node.writeTaskId(taskId, conn)

	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0666)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Errorf("failed to open file, path=%v, err=%v", filePath, err)
		return
	}
	defer file.Close()

	log.Infof("start downloading file: %s", filePath)
	written, err := io.Copy(file, conn)
	if err != nil {
		os.RemoveAll(filePath)
		log.Errorf("failed to download file %s, err=%v", filePath, err)
		return
	}
	atomic.AddInt64(&node.Usage, written)
	log.Infof("download done: %s", filePath)
}

func (node *StorageNode) UploadFile(taskId types.TaskId, filePath, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("failed to connect to the peer, addr=%v, err=%v", addr, err)
		return
	}
	defer conn.Close()

	node.writeTaskId(taskId, conn)

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("failed to open file, path=%v, err=%v", filePath, err)
		return
	}
	defer file.Close()

	log.Infof("start uploading file: %s", filePath)
	io.Copy(conn, file)
	log.Infof("upload done: %s", filePath)
}

func (node *StorageNode) writeTaskId(taskId types.TaskId, conn net.Conn) {
	binary.Write(conn, binary.LittleEndian, taskId)
}
