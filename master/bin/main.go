package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CyDrive/config"
	. "github.com/CyDrive/master"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	conf config.Config
)

func init() {
	logFile, err := os.OpenFile(fmt.Sprintf("master: %v.log", utils.GetDateTimeNow()), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetReportCaller(true)

	configBytes, err := ioutil.ReadFile("master-config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(configBytes, &conf); err != nil {
		panic(err)
	}
	log.Infof("config: %+v", conf)

	config.IpAddr = conf.Ipv4Addr
}

func main() {
	master := NewMaster(conf)

	log.Info("start master")
	master.Start()
}
