package core

import (
	"LibraryManagementV1/LM_V4/config"
	"LibraryManagementV1/LM_V4/global"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const ConfigFile = "./cfg.yaml"

func InitYaml(configFile string) {
	c := &config.Config{}
	yamlConf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("读取配置文件错误:%+v", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalln("yaml解析文件错误")
	}
	log.Println("读取配置文件成功")
	fmt.Println(c)
	global.Config = c
}
