package service

import (
	"context"
	"errors"

	"github.com/jiojioo/gin_template/internal/repo"
	"github.com/jiojioo/gin_template/pkg/hash"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
	redislib "github.com/redis/go-redis/v9"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type UserService struct {
	users       repo.UserRepository
	redisClient *redislib.Client
	jwtCfg      jwtutil.Config
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

type GetUserInfoResp struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Status   int    `json:"status"`
}

func NewUserService(users repo.UserRepository, redisClient *redislib.Client, jwtCfg jwtutil.Config) *UserService {
	jwtutil.Configure(jwtCfg)
	return &UserService{users: users, redisClient: redisClient, jwtCfg: jwtCfg}
}

func (s *UserService) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !hash.Check(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}
	jwtutil.Configure(s.jwtCfg)
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResp{UserID: user.ID, Token: token}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, userID uint64) (*GetUserInfoResp, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetUserInfoResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Status:   user.Status,
	}, nil
}
