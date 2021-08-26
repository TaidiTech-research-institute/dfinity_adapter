package main

import (
	"dfinity_adapter/adapter"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Conf struct {
	Server Server `yaml:"server"`
}
type Server struct {
	Endpoint  string `yaml:"endpoint"`
	LocalPort string `yaml:"local_port"`
	Scrkey    string `yaml:"pri_key"`
	Pubkey    string `yaml:"pub_key"`
}



func main() {
	fmt.Println("Starting dfinity adapter")

	conf := new(Conf)
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Println(err)
	}

	adapterClient, err := adapter.NewdfinityAdaptor(conf.Server.Endpoint, conf.Server.Scrkey, conf.Server.Pubkey)
	if err != nil {
		panic(err)
	}
	adapter.RunWebserver(adapterClient.Handle, conf.Server.LocalPort)
}
