package internal

import (
	"github.com/nzmxd/bserver/config"
	"github.com/nzmxd/bserver/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
// Author [SliverHorn](https://github.com/SliverHorn)
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB
	switch global.CONFIG.System.DbType {
	case "mysql":
		general = global.CONFIG.Mysql.GeneralDB
	case "pgsql":
		general = global.CONFIG.Pgsql.GeneralDB
	case "oracle":
		general = global.CONFIG.Oracle.GeneralDB
	case "sqlite":
		general = global.CONFIG.Sqlite.GeneralDB
	case "mssql":
		general = global.CONFIG.Mssql.GeneralDB
	default:
		general = global.CONFIG.Mysql.GeneralDB
	}
	return &gorm.Config{
		Logger: logger.New(NewWriter(general), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
