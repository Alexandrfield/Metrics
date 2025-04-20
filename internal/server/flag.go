package server

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Alexandrfield/Metrics/internal/common"
)

func getKeyFromFile(path string) []byte {
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

func parseFlags(config *Config) error {
	flag.StringVar(&config.ServerAdderess, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.StringVar(&config.FileStoregePath, "f", "localStorage.dat",
		"path to file for save metrics [default:localStorage.dat]")
	flag.StringVar(&config.DatabaseDsn, "d", "",
		"parametrs for connect Postgress databases [default:-]")
	var signKey string
	flag.StringVar(&signKey, "k", "",
		"key for check sign")
	flag.IntVar(&config.StoreIntervalSecond, "i", 300,
		"interval in seconds for save results on disk [default:300]")
	flag.BoolVar(&config.Restore, "r", true,
		"bool param if we need read exists file with  metrics [default:true]")
	var pathToKey string
	flag.StringVar(&pathToKey, "crypto-key", "",
		"path to crypto key")
	flag.Parse()

	if envServerAdderess, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAdderess = envServerAdderess
	}
	if envSignKey, ok := os.LookupEnv("KEY"); ok {
		signKey = envSignKey
	}
	var err error
	config.SignKey, err = common.GetKeyFromString(signKey)
	if err != nil {
		return fmt.Errorf("try get sign key: %w", err)
	}
	if envFileStoregePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoregePath = envFileStoregePath
	}
	if envDatabaseDsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		config.DatabaseDsn = envDatabaseDsn
		config.FileStoregePath = ""
	}
	if envStoreIntervalSecond, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		value, err := strconv.Atoi(envStoreIntervalSecond)
		if err != nil {
			return fmt.Errorf("try atoi STORE_INTERVAL value; err: %w", err)
		} else {
			config.StoreIntervalSecond = value
		}
	}
	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			return fmt.Errorf("try atoi RESTORE value; err: %w", err)
		} else {
			config.Restore = value
		}
	}
	if envPathCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		pathToKey = envPathCryptoKey
	}
	config.CryptoKeySec = getKeyFromFile(pathToKey)
	return nil
}

// GetServerConfig get server config from environment variables anf flags. Env is preference.
func GetServerConfig() (Config, error) {
	var config Config
	err := parseFlags(&config)
	if err != nil {
		return config, fmt.Errorf("GetServerConfig err:%w", err)
	}
	return config, nil
}
