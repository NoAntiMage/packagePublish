package httpServer

import (
	"PackageServer/config"
	"PackageServer/router"
	"fmt"
	"net/http"
)

type Server struct {
	HttpServer *http.Server
}

func InitHttpServer() *http.Server {
	r := router.NewRouter()
	r.RouterInit()

	port := fmt.Sprintf(":%v", config.ServerConf.Port)

	srv := &http.Server{
		Addr:    port,
		Handler: r.Engine,
	}
	return srv
}
