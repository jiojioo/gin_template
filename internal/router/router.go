// Package router assembles Gin routes.
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/config"
	"github.com/jiojioo/gin_template/internal/handler"
	"github.com/jiojioo/gin_template/internal/middleware"
	routerv1 "github.com/jiojioo/gin_template/internal/router/v1"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
)

func InitRouter(cfg *config.Config, handlers *handler.Handler, jwtCfg jwtutil.Config) *gin.Engine {
	if cfg != nil && cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	}
	engine := gin.New()
	engine.Use(middleware.Recovery(), middleware.RequestLogger(), middleware.CORS())
	routerv1.RegisterUserRoutes(engine, handlers, jwtCfg)
	return engine
}
