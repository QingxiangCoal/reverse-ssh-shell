package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func config() (ADDRESS string, PORT string, USERNAME string) {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	var conf Conf
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	ADDRESS = conf.ADDRESS
	PORT = conf.PORT
	USERNAME = conf.USERNAME
	return ADDRESS, PORT, USERNAME
}

type Conf struct {
	ADDRESS  string
	PORT     string
	USERNAME string
}
