package config

import "fmt"

var (
	IpAddr = "123.57.39.79"
)

type Config struct {
	// "rdb" or "mem"
	AccountStoreType string

	// mysql host or json filepath
	AccountStorePath string

	// only for rdb
	User     string
	Password string
	Database string

	EnvType int
}

func (config Config) PackDSN() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.AccountStorePath,
		config.Database)
}
