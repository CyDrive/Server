package config

import "fmt"

var (
	IpAddr = "123.57.39.79"
)

type Config struct {
	// "rdb" or "mem"
	AccountStoreType string

	// only for rdb
	DatabaseAddr     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string

	EnvType int
}

func (config Config) PackDSN() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseAddr,
		config.DatabaseName)
}
