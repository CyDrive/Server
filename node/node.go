package node

import (
	"fmt"
	"os"
	"time"

	"github.com/CyDrive/node/config"
	"github.com/CyDrive/rpc"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type StorageNode struct {
	*config.Config
	Usage int64
	Id    int32

	heartBeatTimer   *time.Timer
	conn             *grpc.ClientConn
	manageClient     rpc.ManageClient
	fileStreamClient rpc.FileStreamClient
}

func NewStorageNode(config config.Config) *StorageNode {
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
		Config: &config,
		Usage:  usage,

		heartBeatTimer: time.NewTimer(250 * time.Millisecond),
	}

	return &node
}

func (node *StorageNode) Start() {
	var err error

	log.Infof("connecting to the gRPC services")
	node.conn, err = grpc.Dial(node.MasterAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	node.manageClient = rpc.NewManageClient(node.conn)
	node.fileStreamClient = rpc.NewFileStreamClient(node.conn)

	log.Infof("connections setups")
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
	for {
		select {
		case <-node.heartBeatTimer.C:
			node.heartBeatTimer.Reset(250 * time.Millisecond)
			node.HeartBeat()
		}
	}
}
