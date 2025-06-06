package global

import (
	"fmt"
	"github.com/nzmxd/bserver/config"
	"github.com/nzmxd/bserver/utils/timer"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB                  *gorm.DB
	DBList              map[string]*gorm.DB
	REDIS               redis.UniversalClient
	REDISList           map[string]redis.UniversalClient
	ROUTERList          []func(*gin.RouterGroup)
	MONGO               *qmgo.QmgoClient
	CONFIG              config.Server
	VP                  *viper.Viper
	LOG                 *zap.Logger
	Timer               timer.Timer = timer.NewTimerTask()
	Concurrency_Control             = &singleflight.Group{}
	ROUTERS             gin.RoutesInfo
	ACTIVE_DBNAME       *string
	BlackCache          local_cache.Cache
	lock                sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}
