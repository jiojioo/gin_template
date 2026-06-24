package hash_test

import (
	"testing"

	"github.com/jiojioo/gin_template/pkg/hash"
)

func TestPasswordRoundTrip(t *testing.T) {
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("Make() error = %v", err)
	}
	if encoded == "secret" {
		t.Fatal("Make() returned plaintext password")
	}
	if !hash.Check("secret", encoded) {
		t.Fatal("Check() rejected matching password")
	}
}

func TestCheckRejectsWrongPassword(t *testing.T) {
	encoded, err := hash.Make("secret")
	if err != nil {
		t.Fatalf("Make() error = %v", err)
	}
	if hash.Check("wrong", encoded) {
		t.Fatal("Check() accepted a non-matching password")
	}
}
