package config

type ServerConfig struct {
	ServerName string `json:"server_name"`
	ETCDConfig `json:"etcd"`
}

type ETCDConfig struct {
	Addrs string `json:"address"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}
