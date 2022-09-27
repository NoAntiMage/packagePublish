package config

type RpcServerConfig struct {
	ip     string
	port   int
	secret string
}

var ControllerServerConf RpcServerConfig

var SchedularServerConf RpcServerConfig
