package config

import "github.com/CyDrive/rpc"

type Config struct {
	MasterAddr string              `json:"master_addr"`
	Cap        int64               `json:"cap"`
	StorePath  string              `json:"store_path"`
	Type       rpc.StorageNodeType `json:"type"`
}
