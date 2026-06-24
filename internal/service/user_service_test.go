package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jiojioo/gin_template/internal/model"
	"github.com/jiojioo/gin_template/internal/repo"
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

func TestUserServiceLoginSignsJWTForValidPassword(t *testing.T) {
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("hash.Make() error = %v", err)
	}
	svc := NewUserService(fakeUsers{user: &model.User{
		BaseModel: model.BaseModel{ID: 42},
		Username:  "alice",
		Password:  encoded,
		Nickname:  "Alice",
		Status:    1,
	}}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	resp, err := svc.Login(context.Background(), &LoginReq{Username: "alice", Password: "secret"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp.UserID != 42 || resp.Token == "" {
		t.Fatalf("Login() response = %#v, want user id and token", resp)
	}

	claims, err := jwtutil.ParseToken(resp.Token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("token user id = %d, want 42", claims.UserID)
	}
}

func TestUserServiceLoginRejectsInvalidPassword(t *testing.T) {
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("hash.Make() error = %v", err)
	}
	svc := NewUserService(fakeUsers{user: &model.User{Password: encoded}}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	if _, err := svc.Login(context.Background(), &LoginReq{Username: "alice", Password: "wrong"}); err == nil {
		t.Fatal("Login() accepted invalid password")
	}
}

func TestUserServiceGetUserInfoReturnsSafeProfile(t *testing.T) {
	svc := NewUserService(fakeUsers{user: &model.User{
		BaseModel: model.BaseModel{ID: 42},
		Username:  "alice",
		Password:  "encoded-password",
		Nickname:  "Alice",
		Status:    1,
	}}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	profile, err := svc.GetUserInfo(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetUserInfo() error = %v", err)
	}
	if profile.ID != 42 || profile.Username != "alice" || profile.Nickname != "Alice" || profile.Status != 1 {
		t.Fatalf("GetUserInfo() = %#v", profile)
	}
}

func TestUserServiceReturnsRepositoryErrors(t *testing.T) {
	failure := errors.New("database down")
	svc := NewUserService(fakeUsers{err: failure}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	if _, err := svc.Login(context.Background(), &LoginReq{Username: "alice", Password: "secret"}); !errors.Is(err, failure) {
		t.Fatalf("Login() error = %v, want %v", err, failure)
	}
	if _, err := svc.GetUserInfo(context.Background(), 42); !errors.Is(err, failure) {
		t.Fatalf("GetUserInfo() error = %v, want %v", err, failure)
	}
}

func TestUserServiceLoginTreatsMissingUserAsInvalidCredentials(t *testing.T) {
	svc := NewUserService(fakeUsers{err: repo.ErrNotFound}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	_, err := svc.Login(context.Background(), &LoginReq{Username: "ghost", Password: "secret"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestUserServiceGetUserInfoMapsMissingUser(t *testing.T) {
	svc := NewUserService(fakeUsers{err: repo.ErrNotFound}, nil, jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	_, err := svc.GetUserInfo(context.Background(), 42)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserInfo() error = %v, want ErrUserNotFound", err)
	}
}
