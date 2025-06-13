package core

import (
	"fmt"
	"github.com/nzmxd/bserver/global"
	"github.com/nzmxd/bserver/initialize"
	"github.com/nzmxd/bserver/utils/upload"
	"go.uber.org/zap"
	"time"
)

func RunServer() {
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	fmt.Printf(`
	欢迎使用 bserver
	当前版本:v1.0.0
	默认前端文件运行地址:http://127.0.0.1:8080
	默认前端文件运行地址:http://127.0.0.1%s
`, address)
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}

// InitializeSystem 初始化系统所有组件
// 提取为单独函数以便于系统重载时调用
func InitializeSystem() {
	global.VP = Viper() // 初始化Viper
	global.LOG = Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.LOG)
	upload.NewOss()
	global.DB = initialize.Gorm() // gorm连接数据库
	global.CH = initialize.GormClickHouse()
	initialize.DBList()
	if global.CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		if global.CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}
	initialize.SetupHandlers() // 注册全局函数
}
