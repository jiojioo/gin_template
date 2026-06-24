package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/handler"
	"github.com/jiojioo/gin_template/internal/middleware"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
)

func RegisterUserRoutes(engine *gin.Engine, handlers *handler.Handler, jwtCfg jwtutil.Config) {
	group := engine.Group("/api/v1/user")
	group.POST("/login", handlers.User.Login)
	group.GET("/info", middleware.Auth(jwtCfg), handlers.User.GetUserInfo)
}
