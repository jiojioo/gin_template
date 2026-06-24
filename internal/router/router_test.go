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

func TestRouterRejectsLoginWithMissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtCfg := jwtutil.Config{Secret: "test-secret", Expire: time.Hour}
	users := fakeUsers{user: &model.User{Username: "alice", Password: "encoded"}}
	services := service.NewService(&repo.Repository{User: users}, nil, jwtCfg)
	engine := router.InitRouter(&config.Config{}, handler.NewHandler(services), jwtCfg)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("missing fields status = %d, want 400; body = %s", resp.Code, resp.Body.String())
	}
}

func TestRouterRejectsLoginWithWhitespaceFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtCfg := jwtutil.Config{Secret: "test-secret", Expire: time.Hour}
	users := fakeUsers{user: &model.User{Username: "alice", Password: "encoded"}}
	services := service.NewService(&repo.Repository{User: users}, nil, jwtCfg)
	engine := router.InitRouter(&config.Config{}, handler.NewHandler(services), jwtCfg)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{"username":"   ","password":"  "}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("whitespace fields status = %d, want 400; body = %s", resp.Code, resp.Body.String())
	}
}

func TestRouterLoginInvalidCredentialsReturnsGenericMessage(t *testing.T) {
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
	}}
	services := service.NewService(&repo.Repository{User: users}, nil, jwtCfg)
	engine := router.InitRouter(&config.Config{}, handler.NewHandler(services), jwtCfg)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{"username":"alice","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("invalid creds status = %d, want 401; body = %s", resp.Code, resp.Body.String())
	}
	if strings.Contains(resp.Body.String(), "record not found") {
		t.Fatalf("invalid creds leaked persistence detail: %s", resp.Body.String())
	}
	var got struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Message != "invalid username or password" {
		t.Fatalf("message = %q, want %q", got.Message, "invalid username or password")
	}
}

func TestRouterUserInfoMissingUserReturnsGenericFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtCfg := jwtutil.Config{Secret: "test-secret", Expire: time.Hour}
	users := fakeUsers{err: repo.ErrNotFound}
	services := service.NewService(&repo.Repository{User: users}, nil, jwtCfg)
	engine := router.InitRouter(&config.Config{}, handler.NewHandler(services), jwtCfg)

	jwtutil.Configure(jwtCfg)
	token, err := jwtutil.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("missing user status = %d, want 404; body = %s", resp.Code, resp.Body.String())
	}
	if strings.Contains(resp.Body.String(), "record not found") || strings.Contains(resp.Body.String(), "gorm") {
		t.Fatalf("missing user leaked persistence detail: %s", resp.Body.String())
	}
	var got struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Message != "user not found" {
		t.Fatalf("message = %q, want %q", got.Message, "user not found")
	}
}
