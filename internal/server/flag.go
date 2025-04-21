package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

func parseFlags(config *Config) error {
	var serverAdderessTemp string
	var fileStoregePathTemp string
	var databaseDsnTemp string
	var signKeyTemp string
	var storeIntervalSecondTemp int
	var pathToKeyTemp string
	var pathToConfigFile string
	flag.StringVar(&pathToConfigFile, "config", "",
		"path to config file")
	flag.StringVar(&serverAdderessTemp, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.StringVar(&fileStoregePathTemp, "f", "localStorage.dat",
		"path to file for save metrics [default:localStorage.dat]")
	flag.StringVar(&databaseDsnTemp, "d", "",
		"parametrs for connect Postgress databases [default:-]")
	flag.StringVar(&signKeyTemp, "k", "",
		"key for check sign")
	flag.IntVar(&storeIntervalSecondTemp, "i", 300,
		"interval in seconds for save results on disk [default:300]")
	flag.BoolVar(&config.Restore, "r", true,
		"bool param if we need read exists file with  metrics [default:true]")
	flag.StringVar(&pathToKeyTemp, "crypto-key", "",
		"path to crypto key")
	flag.Parse()

	if envPathToConfigFile, ok := os.LookupEnv("CONFIG"); ok {
		pathToConfigFile = envPathToConfigFile
	}

	data := common.GetDataFromFile(pathToConfigFile)
	configT := parseJson(data)
	if configT != nil {
		config = configT
	}
	if serverAdderessTemp != "" {
		config.ServerAdderess = serverAdderessTemp
	}
	if fileStoregePathTemp != "" {
		config.FileStoregePath = fileStoregePathTemp
	}
	if databaseDsnTemp != "" {
		config.DatabaseDsn = databaseDsnTemp
	}
	if storeIntervalSecondTemp != 0 {
		config.StoreIntervalSecond = storeIntervalSecondTemp
	}

	if envServerAdderess, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAdderess = envServerAdderess
	}
	if envSignKey, ok := os.LookupEnv("KEY"); ok {
		signKeyTemp = envSignKey
	}
	var err error
	config.SignKey, err = common.GetKeyFromString(signKeyTemp)
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
		pathToKeyTemp = envPathCryptoKey
	}
	if pathToKeyTemp != "" {
		config.CryptoKeySec = common.GetDataFromFile(pathToKeyTemp)
	}
	return nil
}

func parseJson(data []byte) *Config {
	var conf configJSON
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return nil
	}
	var config Config
	config.CryptoKeySec = common.GetDataFromFile(conf.CryptoKeySec)
	config.DatabaseDsn = conf.DatabaseDsn
	config.FileStoregePath = conf.FileStoregePath
	config.Restore = conf.Restore
	config.ServerAdderess = conf.ServerAdderess
	config.SignKey = []byte(conf.SignKey)
	duration, err := time.ParseDuration(conf.StoreIntervalSecond)
	if err != nil {
		return nil
	}
	config.StoreIntervalSecond = int(duration.Seconds())
	return &config
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
