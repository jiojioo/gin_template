package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/config"
	"github.com/jiojioo/gin_template/internal/db"
	"github.com/jiojioo/gin_template/internal/db/mysql"
	redisclient "github.com/jiojioo/gin_template/internal/db/redis"
	"github.com/jiojioo/gin_template/internal/handler"
	"github.com/jiojioo/gin_template/internal/repo"
	"github.com/jiojioo/gin_template/internal/router"
	"github.com/jiojioo/gin_template/internal/service"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
	"github.com/jiojioo/gin_template/pkg/logger"
)

func main() {
	cfg := config.MustLoad("etc/config.yaml")
	if err := logger.Init(cfg.Log); err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()
	if err := logger.InitGinWriter(cfg.Log); err != nil {
		panic(err)
	}
	if err := db.Init(cfg); err != nil {
		panic(err)
	}

	jwtCfg := jwtConfig(cfg)
	repos := repo.NewRepository(mysql.Client)
	services := service.NewService(repos, redisclient.Client, jwtCfg)
	engine := buildRouter(cfg, services, jwtCfg)
	if err := engine.Run(cfg.Server.Addr); err != nil {
		panic(err)
	}
}

func buildRouter(cfg *config.Config, services *service.Service, jwtCfg jwtutil.Config) *gin.Engine {
	return router.InitRouter(cfg, handler.NewHandler(services), jwtCfg)
}

func jwtConfig(cfg *config.Config) jwtutil.Config {
	return jwtutil.Config{
		Secret: cfg.JWT.Secret,
		Expire: time.Duration(cfg.JWT.Expire) * time.Second,
	}
}
