package config

type RedisConfig struct {
	Ip       string
	Port     int
	Password string
}

var RedisConf RedisConfig
