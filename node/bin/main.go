package main

import (
	"github.com/CyDrive/node"
	"github.com/CyDrive/node/config"
	"github.com/CyDrive/rpc"
)

func main() {
	config := config.Config{
		MasterAddr: "123.57.39.79:6455",
		Cap:        5,
		StorePath:  "./data",
		Type:       rpc.StorageNodeType_Private,
	}
	node := node.NewStorageNode(config)

	node.Start()
}
