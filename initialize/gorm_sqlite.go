package initialize

import (
	"github.com/glebarez/sqlite"
	"github.com/nzmxd/bserver/config"
	"github.com/nzmxd/bserver/global"
	"github.com/nzmxd/bserver/initialize/internal"
	"gorm.io/gorm"
)

// GormSqlite 初始化Sqlite数据库
func GormSqlite() *gorm.DB {
	s := global.CONFIG.Sqlite
	return initSqliteDatabase(s)
}

// GormSqliteByConfig 初始化Sqlite数据库用过传入配置
func GormSqliteByConfig(s config.Sqlite) *gorm.DB {
	return initSqliteDatabase(s)
}

// initSqliteDatabase 初始化Sqlite数据库辅助函数
func initSqliteDatabase(s config.Sqlite) *gorm.DB {
	if s.Dbname == "" {
		return nil
	}

	if db, err := gorm.Open(sqlite.Open(s.Dsn()), internal.Gorm.Config(s.Prefix, s.Singular)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}
