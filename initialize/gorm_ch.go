package initialize

import (
	"git.zingfront.cn/liubin/bserver/global"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func GormClickHouse() *gorm.DB {
	m := global.CONFIG.ClickHouse
	if m.Dbname == "" {
		return nil
	}
	clickhouseConfig := clickhouse.Config{
		DSN: m.Dsn(),
	}
	if db, err := gorm.Open(clickhouse.New(clickhouseConfig), &gorm.Config{
		QueryFields: true,
	}); err != nil {
		panic(err)
	} else {
		return db
	}
}
