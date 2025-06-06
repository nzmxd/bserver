package initialize

import (
	"github.com/nzmxd/bserver/global"
)

func bizModel() error {
	db := global.DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
