package config

import "github.com/CyDrive/rpc"

type Config struct {
	MasterAddr string
	Cap        int64
	StorePath  string
	Type       rpc.StorageNodeType
}
