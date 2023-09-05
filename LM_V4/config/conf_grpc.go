package config

import (
	"fmt"
	"strconv"
)

type GRPC struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (r GRPC) Addr() string {
	fmt.Printf("Config GRPC :%+v\n", r.Host+":"+strconv.Itoa(r.Port))
	return r.Host + ":" + strconv.Itoa(r.Port)
}
