package main

import (
	"fastProxy/app/config"
	"fastProxy/app/core"
	"fmt"
)

func main() {

	fmt.Printf("程序已运行至：%d 端口\n", config.GlobalConfig.Server.Port)
	core.Start()
}
