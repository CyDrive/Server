package config

import "github.com/yah01/CyDrive/rpc"

type Config struct {
	MasterAddr string
	Cap        int64
	StorePath  string
	Type       rpc.StorageNodeType
}
