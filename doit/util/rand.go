package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

func RandString() (string, error) {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, err = hash.Write(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
