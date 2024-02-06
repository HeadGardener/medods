package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 8
)

func GetStringHash(str string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(str), cost)

	return string(passwordHash)
}

func CompareHashAndString(passHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))

	return err == nil
}
