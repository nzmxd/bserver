package core

import (
	"git.zingfront.cn/liubin/bserver/global"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// TimerTask 包装任务函数，处理 panic 和错误日志
func TimerTask(f func() error) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				global.LOG.Error("timer task panic recovered",
					zap.Any("recover", r),
					zap.Stack("stack"),
				)
			}
		}()

		if f == nil {
			global.LOG.Warn("timer task is nil, skipping execution")
			return
		}

		if err := f(); err != nil {
			global.LOG.Error("timer task execution failed", zap.Error(err))
		}
	}
}

// RegisterCronTask 注册单个定时任务
func RegisterCronTask(name, spec string, taskFunc func() error, description string, options ...cron.Option) {
	if taskFunc == nil {
		global.LOG.Warn("cron task function is nil, skipping registration", zap.String("task", name))
		return
	}

	opts := append([]cron.Option{cron.WithSeconds()}, options...)
	_, err := global.Timer.AddTaskByFunc(name, spec, TimerTask(taskFunc), description, opts...)
	if err != nil {
		global.LOG.Error("failed to register cron task",
			zap.String("task", name),
			zap.String("spec", spec),
			zap.String("desc", description),
			zap.Error(err),
		)
		return
	}

	global.LOG.Info("cron task registered successfully",
		zap.String("task", name),
		zap.String("spec", spec),
		zap.String("desc", description),
	)
}
