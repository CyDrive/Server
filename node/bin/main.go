package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/CyDrive/node"
	"github.com/CyDrive/node/config"
)

func main() {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var config config.Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	node := node.NewStorageNode(&config)

	node.Start()
}
