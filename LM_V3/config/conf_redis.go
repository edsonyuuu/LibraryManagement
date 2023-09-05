package config

import (
	"fmt"
	"strconv"
)

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

func (r Redis) Addr() string {
	fmt.Printf("Config Redis :%+v\n", r.Host+":"+strconv.Itoa(r.Port))
	return r.Host + ":" + strconv.Itoa(r.Port)
}
