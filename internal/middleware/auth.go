// Package middleware contains Gin middleware.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
	"github.com/jiojioo/gin_template/pkg/response"
)

func Auth(jwtCfg jwtutil.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}
		jwtutil.Configure(jwtCfg)
		claims, err := jwtutil.ParseToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "invalid bearer token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
