package main

import (
	"git.zingfront.cn/liubin/bserver/core"
)

func main() {
	// 初始化系统
	core.InitializeSystem()
	// 运行服务器
	core.RunServer()
}
