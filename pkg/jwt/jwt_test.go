package jwt_test

import (
	"testing"
	"time"

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
