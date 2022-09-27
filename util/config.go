package util

import (
	"fmt"

	"PackageServer/config"

	"github.com/spf13/viper"
)

var AppConfig *viper.Viper

func ConfigInit(serverName string) {
	configName := fmt.Sprintf("./conf/%v.yaml", serverName)

	AppConfig = viper.New()
	AppConfig.SetConfigFile(configName)
	if err := AppConfig.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Config init...")
	config.ServerConf.Name = AppConfig.Get("name").(string)
	config.ServerConf.AreaName = AppConfig.Get("area").(string)
	config.ServerConf.Port = AppConfig.Get("port").(int)
	config.ServerConf.RouterMode = AppConfig.Get("routerMode").(string)
	config.ServerConf.LoggerMode = AppConfig.Get("loggerMode").(string)
	config.ServerConf.DataDir = AppConfig.Get("dataDir").(string)
	config.ServerConf.PackageDir = AppConfig.Get("packageDir").(string)
	fmt.Println(config.ServerConf)

	redisConf := AppConfig.Get("redis").(map[string]interface{})
	config.RedisConf.Ip = redisConf["ip"].(string)
	config.RedisConf.Port = redisConf["port"].(int)
	config.RedisConf.Password = redisConf["password"].(string)

}
