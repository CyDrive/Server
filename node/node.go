package node

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yah01/CyDrive/node/config"
	"github.com/yah01/CyDrive/rpc"
	"github.com/yah01/CyDrive/utils"
	"google.golang.org/grpc"
)

type StorageNode struct {
	*config.Config
	Usage int64
	Id    int32

	heartBeatTimer *time.Timer
	conn       *grpc.ClientConn
	grpcClient rpc.ManageClient

	log *logrus.Logger
}

func NewStorageNode(config config.Config) *StorageNode {
	usage, err := utils.DirSize(config.StorePath)
	if err != nil {
		panic(err)
	}

	logfile, err := os.Create(time.Now().String() + ".log")
	if err != nil {
		panic(err)
	}

	log := logrus.New()
	log.Out = logfile

	node := StorageNode{
		Config: &config,
		Usage:  usage,

		heartBeatTimer: time.NewTimer(250 * time.Millisecond),

		log:    log,
	}

	return &node
}

func (node *StorageNode) Start() {
	var err error
	node.conn, err = grpc.Dial(node.MasterAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	node.grpcClient = rpc.NewManageClient(node.conn)

	for {
		if err := node.JoinCluster(); err != nil {
			node.log.Warn(err)
			continue
		}
		break
	}

	go func ()  {
		for {
			select {
			case <-node.heartBeatTimer.C:
				node.heartBeatTimer.Reset(250 * time.Millisecond)
				node.HeartBeat()
			}
		}
	}
}

func (node *StorageNode) JoinCluster() error {
	req := &rpc.JoinClusterRequest{
		Capacity: node.Cap,
		Usage:    node.Usage,
		Type:     node.Type,
	}

	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	resp, err := node.grpcClient.JoinCluster(ctx, req)
	if err != nil {
		return err
	}

	node.Id = resp.Id

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
	_, err := node.grpcClient.HeartBeats(ctx, req)
	if err != nil {
		node.log.Warn(err)
	}
}
