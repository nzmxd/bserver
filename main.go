package main

import (
	"github.com/nzmxd/bserver/core"
	"github.com/nzmxd/bserver/global"
	"github.com/nzmxd/bserver/initialize"
	"go.uber.org/zap"
)

func main() {
	// 初始化系统
	initializeSystem()
	// 运行服务器
	core.RunServer()
}

// initializeSystem 初始化系统所有组件
// 提取为单独函数以便于系统重载时调用
func initializeSystem() {
	global.VP = core.Viper() // 初始化Viper
	global.LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.LOG)
	global.DB = initialize.Gorm() // gorm连接数据库
	initialize.DBList()
	initialize.SetupHandlers() // 注册全局函数
}
