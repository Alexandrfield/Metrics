package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GetKeyFromString(key string) ([]byte, error) {
	data, err := hex.DecodeString(key)
	if err != nil {
		return []byte{}, fmt.Errorf("Sign error:%w", err)
	}
	return data, nil
}
func Sign(msg []byte, signKey []byte) ([]byte, error) {
	// data, err := hex.DecodeString(msg)
	// if err != nil {
	// 	return []byte{}, fmt.Errorf("Sign error:%w", err)
	// }
	h := hmac.New(sha256.New, signKey)
	h.Write(msg)
	return h.Sum(nil), nil
}
