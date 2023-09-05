package config

type Config struct {
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
	Logger Logger `json:"logger" yaml:"logger"`
	System System `json:"system" yaml:"system"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
	Redis  Redis  `yaml:"redis"`
	SMS    SMS    `yaml:"sms"`
	GRPC   GRPC   `yaml:"grpc"`
}
