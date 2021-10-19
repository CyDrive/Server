package main

import (
	"github.com/CyDrive/node"
	"github.com/CyDrive/node/config"
	"github.com/CyDrive/rpc"
)

func main() {
	config := config.Config{
		MasterAddr: "127.0.0.1:6455",
		Cap:        5,
		StorePath:  ".",
		Type:       rpc.StorageNodeType_Private,
	}
	node := node.NewStorageNode(config)

	node.Start()
}
