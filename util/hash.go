package util

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var hashCost, _ = strconv.Atoi(os.Getenv("HASH_COST"))

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(hash), err
}

func VerifyPassword(hash, password string) bool {
	err :=bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
