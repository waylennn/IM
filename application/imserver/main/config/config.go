package config

import (
	"github.com/micro/go-micro/v2/config"
)

type ServerConfig struct {
	ServerName string `json:"server_name"`
	ImServerList []*ImRpc `json:"im_server_list"`
	ETCDConfig `json:"etcd"`
	RateLimit  `json:"rate_limit"`
	DB         `json:"db"`
	RabbitMq   `json:"rabbitMq"`

}

type (
	ETCDConfig struct {
		Addrs string `json:"address"`
	}

	RateLimit struct {
		Rate int `json:"rate"`
	}

	DB struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	}

	RabbitMq struct {
		Address []string `json:"address"`
		Topic   string   `json:"topic"`
	}

	ImRpc struct {
		Address     string `json:"address"`// 这是一个真正 的ip地址,连接IM服务的地址
		AmqbAddress []string `json:"amqb_address"`
		Topic       string `json:"topic"`
		ServerName  string `json:"server_name"`
		Weight int `json:"weight"`
	}
)

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func InitConfig(defaultPath string) (conf *ServerConfig, err error) {
	//初始化配置文件
	conf = NewServerConfig()
	if err := config.LoadFile(defaultPath); err != nil {
		return nil, err
	}
	if err := config.Scan(conf); err != nil {
		return nil, err
	}
	return
}
