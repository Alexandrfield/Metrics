package common

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int, logger Loger) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logger.Errorf("rsa generate key. err:%s", err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("marshal pki.err:%f", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt pem blockerr:%f", err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, fmt.Errorf("parse pks blockerr:%f", err)
	}
	return key, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt pem block. err:%f", err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, fmt.Errorf("parse pub key. err:%f", err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key not ok")
	}
	return key, nil
}
func EncryptData(msg []byte, key []byte, logger Loger) []byte {
	var res []byte = msg
	if len(msg) > 0 && len(key) > 0 {
		pybKey, err := BytesToPublicKey(key)
		if err != nil {
			logger.Errorf("problem with rsa key")
			return res
		}
		encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pybKey, msg, nil)
		if err != nil {
			logger.Errorf("error encrypting message with public key: %w", err)
			return res
		}
		res = encryptedData
	}
	return res
}
func DecryptData(msg []byte, key []byte, logger Loger) []byte {
	var res []byte = msg
	if len(msg) > 0 && len(key) > 0 {
		privKey, err := BytesToPrivateKey(key)
		if err != nil {
			logger.Errorf("problem with rsa key")
			return res
		}
		plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, msg, nil)
		if err != nil {
			logger.Errorf("error encrypting message with public key: %w", err)
			return res
		}
		res = plaintext
	}
	return res
}
