// Package service implements business use cases.
package service

import (
	"github.com/jiojioo/gin_template/internal/repo"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
	redislib "github.com/redis/go-redis/v9"
)

type Service struct {
	User *UserService
}

func NewService(repos *repo.Repository, redisClient *redislib.Client, jwtCfg jwtutil.Config) *Service {
	return &Service{User: NewUserService(repos.User, redisClient, jwtCfg)}
}
