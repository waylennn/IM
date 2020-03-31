package user_config

import "github.com/micro/go-micro/v2/config"

type ServerConfig struct {
	ServerName string `json:"server_name"`
	ETCDConfig `json:"etcd"`
	RateLimit  `json:"rate_limit"`
	DB `json:"db"`
}

type ETCDConfig struct {
	Addrs string `json:"address"`
}

type RateLimit struct {
	Rate int `json:"rate"`
}

type DB struct{
	Address string `json:"address"`
	Name string `json:"name"`
}
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func InitConfig() (conf *ServerConfig, err error) {
	//初始化配置文件
	conf = NewServerConfig()
	if err := config.LoadFile("./userconfig/config.json"); err != nil {
		return nil, err
	}
	if err := config.Scan(conf); err != nil {
		return nil, err
	}
	return
}
