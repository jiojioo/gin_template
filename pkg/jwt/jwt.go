// Package jwt creates and validates user access tokens.
package jwt

import (
	"errors"
	"sync"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret string
	Expire time.Duration
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwtlib.RegisteredClaims
}

var (
	configMu sync.RWMutex
	config   Config
)

func Configure(next Config) {
	configMu.Lock()
	defer configMu.Unlock()
	config = next
}

func GenerateToken(userID uint64) (string, error) {
	cfg := currentConfig()
	if cfg.Secret == "" {
		return "", errors.New("jwt secret is not configured")
	}
	if cfg.Expire <= 0 {
		return "", errors.New("jwt expiration must be positive")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(cfg.Expire)),
		},
	}
	return jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte(cfg.Secret))
}

func ParseToken(raw string) (*Claims, error) {
	cfg := currentConfig()
	if cfg.Secret == "" {
		return nil, errors.New("jwt secret is not configured")
	}

	claims := new(Claims)
	token, err := jwtlib.ParseWithClaims(raw, claims, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok || token.Method.Alg() != jwtlib.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected jwt signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid jwt token")
	}
	if claims.ExpiresAt == nil {
		return nil, errors.New("jwt token expiration is required")
	}
	return claims, nil
}

func currentConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}
