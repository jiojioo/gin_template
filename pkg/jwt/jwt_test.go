package jwt_test

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	jwtutil "github.com/jiojioo/gin_template/pkg/jwt"
)

func TestGenerateAndParseToken(t *testing.T) {
	jwtutil.Configure(jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	token, err := jwtutil.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	claims, err := jwtutil.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("UserID = %d, want 42", claims.UserID)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	jwtutil.Configure(jwtutil.Config{Secret: "signing-secret", Expire: time.Hour})
	token, err := jwtutil.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	jwtutil.Configure(jwtutil.Config{Secret: "other-secret", Expire: time.Hour})

	if _, err := jwtutil.ParseToken(token); err == nil {
		t.Fatal("ParseToken() accepted token signed with a different secret")
	}
}

func TestParseTokenRejectsTokenWithoutExpiration(t *testing.T) {
	jwtutil.Configure(jwtutil.Config{Secret: "test-secret", Expire: time.Hour})

	claims := jwtutil.Claims{
		UserID: 42,
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt: jwtlib.NewNumericDate(time.Now()),
		},
	}
	token, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := jwtutil.ParseToken(token); err == nil {
		t.Fatal("ParseToken() accepted token without expiration")
	}
}
