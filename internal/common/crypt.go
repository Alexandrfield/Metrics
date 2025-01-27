package common

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GetKeyFromString(key string) ([]byte, error) {
	return []byte(key), nil
}
func Sign(msg []byte, signKey []byte) ([]byte, error) {
	h := hmac.New(sha256.New, signKey)
	h.Write(msg)
	return h.Sum(nil), nil
}

func CheckHash(msg []byte, msgSign []byte, signKey []byte) bool {
	actualSign, _ := Sign(msg, signKey)
	return hmac.Equal(actualSign, msgSign)
}
