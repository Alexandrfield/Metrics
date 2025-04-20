package agent

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

func parseFlags(conf *Config) error {
	flag.StringVar(&conf.ServerAdderess, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.IntVar(&conf.ReportIntervalSecond, "r", 10,
		"interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&conf.PollIntervalSecond, "p", 2,
		"interval in seconds for check metrics [default: 2 second]")
	flag.IntVar(&conf.RateLimit, "l", 1,
		"limit count reqyst in time [default: 1]")
	var signKey string
	flag.StringVar(&signKey, "k", "",
		"key for sign [default: nil]")
	var pathToKey string
	flag.StringVar(&pathToKey, "crypto-key", "",
		"path to file with key for [default: nil]")
	flag.Parse()

	if envServerAdderess := os.Getenv("ADDRESS"); envServerAdderess != "" {
		conf.ServerAdderess = envServerAdderess
	}

	if envSignKey := os.Getenv("KEY"); envSignKey != "" {
		signKey = envSignKey
	}
	var err error
	conf.SignKey, err = common.GetKeyFromString(signKey)
	if err != nil {
		return fmt.Errorf("try get sign key: %w", err)
	}

	if envRateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		value, err := strconv.Atoi(envRateLimit)
		if err != nil {
			return fmt.Errorf("try atoi . value; err: %w", err)
		} else {
			conf.RateLimit = value
		}
	}
	if envReportIntervalSecond, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		value, err := strconv.Atoi(envReportIntervalSecond)
		if err != nil {
			return fmt.Errorf("try atoi REPORT_INTERVAL value; err: %w", err)
		} else {
			conf.ReportIntervalSecond = value
		}
	}
	if envPollIntervalSecond, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		value, err := strconv.Atoi(envPollIntervalSecond)
		if err != nil {
			return fmt.Errorf("try atoi POLL_INTERVAL value; err: %w", err)
		} else {
			conf.PollIntervalSecond = value
		}
	}
	if envPathCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		pathToKey = envPathCryptoKey
	}
	conf.CryptoKeyOpen = getKeyFromFile(pathToKey)
	return nil
}

// GetAgentConfig get server config from environment variables anf flags. Env is preference.
func GetAgentConfig() (Config, error) {
	var config Config
	err := parseFlags(&config)
	if err != nil {
		return config, fmt.Errorf("GetAgentConfig err:%w", err)
	}
	return config, nil
}
