package router_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/config"
	"github.com/jiojioo/gin_template/internal/handler"
	"github.com/jiojioo/gin_template/internal/model"
	"github.com/jiojioo/gin_template/internal/repo"
	"github.com/jiojioo/gin_template/internal/router"
	"github.com/jiojioo/gin_template/internal/service"
	"github.com/jiojioo/gin_template/pkg/hash"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
)

type fakeUsers struct {
	user *model.User
	err  error
}

func (f fakeUsers) FindByID(context.Context, uint64) (*model.User, error) {
	return f.user, f.err
}

func (f fakeUsers) FindByUsername(context.Context, string) (*model.User, error) {
	return f.user, f.err
}

func (f fakeUsers) Create(context.Context, *model.User) error {
	return f.err
}

func TestRouterRegistersLoginAndProtectedUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("hash.Make() error = %v", err)
	}
	jwtCfg := jwtutil.Config{Secret: "test-secret", Expire: time.Hour}
	users := fakeUsers{user: &model.User{
		BaseModel: model.BaseModel{ID: 42},
		Username:  "alice",
		Password:  encoded,
		Nickname:  "Alice",
		Status:    1,
	}}
	services := service.NewService(&repo.Repository{User: users}, nil, jwtCfg)
	engine := router.InitRouter(&config.Config{}, handler.NewHandler(services), jwtCfg)

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{"username":"alice","password":"secret"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	engine.ServeHTTP(loginResp, loginReq)
	if loginResp.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", loginResp.Code, loginResp.Body.String())
	}
	var loginBody struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(loginResp.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login body: %v", err)
	}
	if loginBody.Data.Token == "" {
		t.Fatalf("login token is empty: %s", loginResp.Body.String())
	}

	infoReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	infoReq.Header.Set("Authorization", "Bearer "+loginBody.Data.Token)
	infoResp := httptest.NewRecorder()
	engine.ServeHTTP(infoResp, infoReq)
	if infoResp.Code != http.StatusOK {
		t.Fatalf("info status = %d, body = %s", infoResp.Code, infoResp.Body.String())
	}
	if strings.Contains(infoResp.Body.String(), "encoded") {
		t.Fatalf("user info leaked password hash: %s", infoResp.Body.String())
	}
}
