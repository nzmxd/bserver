package core

import (
	"fmt"
	"github.com/nzmxd/bserver/global"
	"github.com/nzmxd/bserver/initialize"
	"time"
)

func RunServer() {
	if global.CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		if global.CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	fmt.Printf(`
	欢迎使用 bserver
	当前版本:v1.0.0
	默认前端文件运行地址:http://127.0.0.1:8080
	默认前端文件运行地址:http://127.0.0.1:%d
`, address)
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
