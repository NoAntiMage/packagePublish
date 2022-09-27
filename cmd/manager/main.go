package main

import (
	"PackageServer/server"
	"PackageServer/util"
	"fmt"
)

const serverName = "manager"

func main() {
	util.ConfigInit(serverName)
	fmt.Printf("%s Server starting...", serverName)
	server.Start()
}
