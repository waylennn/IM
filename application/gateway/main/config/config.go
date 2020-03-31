package config

import "github.com/micro/go-micro/v2/config"

type ServerConfig struct {
	ServerName string `json:"server_name"`
	ETCDConfig `json:"etcd"`
	RateLimit  `json:"rate_limit"`
	Web        `json:"web"`
	DB         `json:"db"`
	ImRpcServer    `json:"im_rpc_server"`
	UserRpcServer `json:"user_rpc_server"`
}

type (
	ETCDConfig struct {
		Addrs string `json:"address"`
	}

	RateLimit struct {
		Rate int `json:"rate"`
	}

	Web struct {
		Address string `json:"address"`
	}

	DB struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	ImRpcServer   struct {
		ServerName   string `json:"server_name"`
		ImServerList []*ImRpc `json:"im_server_list"`
	}
	ImRpc struct {
		Address     string `json:"address"`// 这是一个真正 的ip地址,连接IM服务的地址
		AmqbAddress []string `json:"amqb_address"`
		Topic       string `json:"topic"`
		ServerName  string `json:"server_name"`
		Weight int `json:"weight"`
	}

	UserRpcServer struct {
		ClientName string `json:"client_name"`
		ServerName string `json:"server_name"`
	}

)

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func InitConfig() (conf *ServerConfig, err error) {
	//初始化配置文件
	conf = NewServerConfig()
	if err := config.LoadFile("./config/config.json"); err != nil {
		return nil, err
	}
	if err := config.Scan(conf); err != nil {
		return nil, err
	}
	return
}
