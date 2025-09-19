package main

import (
	"github.com/nzmxd/bserver/core"
)

func main() {
	// 初始化系统
	core.InitializeSystem()
	// 运行服务器
	core.RunServer()
}
