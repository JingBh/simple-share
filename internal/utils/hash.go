package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func MD5HashBase64(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func HashPassword(password string) (string, error) {
	rawHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawHash), nil
}

func VerifyPassword(password, hash string) error {
	rawHash, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(rawHash, []byte(password))
}
