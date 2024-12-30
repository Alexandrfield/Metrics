package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func parseFlags(config *Config) {
	flag.StringVar(&config.ServerAdderess, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.StringVar(&config.FileStoregePath, "f", "localStorage.dat",
		"path to file for save metrics [default:localStorage.dat]")
	flag.IntVar(&config.StoreIntervalSecond, "i", 300,
		"interval in seconds for save results on disk [default:300]")
	flag.BoolVar(&config.Restore, "r", true,
		"bool param if we need read exists file with  metrics [default:true]")
	flag.Parse()

	if envServerAdderess, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAdderess = envServerAdderess
	}
	if envFileStoregePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoregePath = envFileStoregePath
	}
	if envStoreIntervalSecond, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		value, err := strconv.Atoi(envStoreIntervalSecond)
		if err != nil {
			fmt.Printf("error atoi STORE_INTERVAL; err: %s\n", err)
		} else {
			config.StoreIntervalSecond = value
		}
	}
	if envRestore, ok := os.LookupEnv("RESTORE"); ok {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			fmt.Printf("error atoi RESTORE; err: %v\n", err)
		} else {
			config.Restore = value
		}
	}
}

func GetServerConfig() Config {
	var config Config
	parseFlags(&config)
	return config
}
