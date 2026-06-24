package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/config"
	"github.com/jiojioo/gin_template/internal/model"
	"github.com/jiojioo/gin_template/internal/repo"
	"github.com/jiojioo/gin_template/internal/service"
	"github.com/jiojioo/gin_template/pkg/hash"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
)

type appFakeUsers struct {
	user *model.User
}

func (f appFakeUsers) FindByID(context.Context, uint64) (*model.User, error) {
	return f.user, nil
}

func (f appFakeUsers) FindByUsername(context.Context, string) (*model.User, error) {
	return f.user, nil
}

func (f appFakeUsers) Create(context.Context, *model.User) error {
	return nil
}

func TestApplicationBuildsRouterAfterDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("hash.Make() error = %v", err)
	}
	jwtCfg := jwtutil.Config{Secret: "test-secret", Expire: time.Hour}
	services := service.NewService(&repo.Repository{User: appFakeUsers{user: &model.User{
		BaseModel: model.BaseModel{ID: 42},
		Username:  "alice",
		Password:  encoded,
		Nickname:  "Alice",
		Status:    1,
	}}}, nil, jwtCfg)

	engine := buildRouter(&config.Config{Server: config.ServerConfig{Mode: gin.TestMode}}, services, jwtCfg)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/user/login", nil)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	if resp.Code == http.StatusNotFound {
		t.Fatal("router did not register user routes")
	}
}
