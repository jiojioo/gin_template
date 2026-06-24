// Package db initializes infrastructure clients required at startup.
package db

import (
	"github.com/jiojioo/gin_template/internal/config"
	"github.com/jiojioo/gin_template/internal/db/mysql"
	"github.com/jiojioo/gin_template/internal/db/redis"
)

var (
	mysqlInit   = mysql.Init
	autoMigrate = mysql.AutoMigrate
	initData    = mysql.InitData
	redisInit   = redis.Init
)

func Init(cfg *config.Config) error {
	if err := mysqlInit(cfg.MySQL); err != nil {
		return err
	}
	if err := autoMigrate(); err != nil {
		return err
	}
	if err := initData(); err != nil {
		return err
	}
	return redisInit(cfg.Redis)
}

func replaceInitHooks(
	nextMySQLInit func(config.MySQLConfig) error,
	nextAutoMigrate func() error,
	nextInitData func() error,
	nextRedisInit func(config.RedisConfig) error,
) func() {
	previousMySQLInit := mysqlInit
	previousAutoMigrate := autoMigrate
	previousInitData := initData
	previousRedisInit := redisInit

	mysqlInit = nextMySQLInit
	autoMigrate = nextAutoMigrate
	initData = nextInitData
	redisInit = nextRedisInit

	return func() {
		mysqlInit = previousMySQLInit
		autoMigrate = previousAutoMigrate
		initData = previousInitData
		redisInit = previousRedisInit
	}
}
