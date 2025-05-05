package common

import (
	"log"
	"os"
)

func GetDataFromFile(path string) []byte {
	if path == "" {
		return []byte{}
	}
	fContent, err := os.ReadFile(path)
	if err != nil {
		log.Printf("promlem read file. err:%s", err)
		return []byte{}
	}
	return fContent
}
