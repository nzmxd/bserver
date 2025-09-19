package initialize

import (
	"git.zingfront.cn/liubin/bserver/global"
)

func bizModel() error {
	db := global.DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
