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
	logFile, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetReportCaller(true)

}

func main() {
	var accountStore store.AccountStore = nil
	if conf.UserStoreType == "mem" {
		accountStore = store.NewMemStore("user_data/user.json")
	}
	master := NewMaster(conf, env.NewLocalEnv(), accountStore)
	master.Start()
}
