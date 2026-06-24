// Package hash provides bcrypt password hashing helpers.
package hash

import "golang.org/x/crypto/bcrypt"

func Make(password string) (string, error) {
	encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func Check(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
