package config

import "fmt"

var (
	IpAddr string = "0.0.0.0"
)

type Config struct {
	Ipv4Addr         string `yaml:"ipv4_addr"`
	AccountStoreType string `yaml:"account_store_type"`

	// only for relational database
	DatabaseAddr     string `yaml:"database_addr,omitempty"`
	DatabaseName     string `yaml:"database_name,omitempty"`
	DatabaseUser     string `yaml:"database_user,omitempty"`
	DatabasePassword string `yaml:"database_password,omitempty"`

	EnvType string `yaml:"env_type"`
}

func (config Config) PackDSN() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseAddr,
		config.DatabaseName)
}
