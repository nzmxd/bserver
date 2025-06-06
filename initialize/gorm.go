package initialize

import (
	"github.com/nzmxd/bserver/global"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.CONFIG.System.DbType {
	case "mysql":
		global.ACTIVE_DBNAME = &global.CONFIG.Mysql.Dbname
		return GormMysql()
	case "sqlite":
		global.ACTIVE_DBNAME = &global.CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.ACTIVE_DBNAME = &global.CONFIG.Mysql.Dbname
		return GormMysql()
	}
}
