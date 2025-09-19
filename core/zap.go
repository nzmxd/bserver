package core

import (
	"fmt"
	"git.zingfront.cn/liubin/bserver/core/internal"
	"git.zingfront.cn/liubin/bserver/global"
	"git.zingfront.cn/liubin/bserver/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.CONFIG.Zap.Director)
		_ = os.Mkdir(global.CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger = zap.New(zapcore.NewTee(cores...))
	if global.CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
