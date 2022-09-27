package config

type ServerConfig struct {
	Name       string
	AreaName   string
	Port       int
	RouterMode string
	LoggerMode string
	DataDir    string
	PackageDir string
}

var ServerConf ServerConfig
