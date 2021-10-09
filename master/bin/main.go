package main

import (
	"os"

	"github.com/CyDrive/config"
	. "github.com/CyDrive/master"
	"github.com/CyDrive/master/env"
	"github.com/CyDrive/master/store"
	log "github.com/sirupsen/logrus"
)

var (
	conf config.Config

	isOnline      bool
	serverAddress string
)

func init() {
	logFile, err := os.OpenFile("master.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetReportCaller(true)

	conf.AccountStoreType = "mem"
	conf.AccountStorePath = "user_data/user.json"
}

func main() {
	var accountStore store.AccountStore = nil
	if conf.AccountStoreType == "mem" {
		accountStore = store.NewMemStore(conf.AccountStorePath)
	}
	master := NewMaster(conf, env.NewLocalEnv(), accountStore)
	master.Start()
}
