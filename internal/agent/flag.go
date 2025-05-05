package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

func parseFlags(config *Config) error {
	var pathToConfigFile string
	var serverAdderessTemp string
	var reportIntervalSecondTemp int
	var pollIntervalSecondTemp int
	var rateLimitTemp int
	var signKeyTemp string
	var pathToKey string
	flag.StringVar(&pathToConfigFile, "config", "",
		"path to config file")
	flag.StringVar(&serverAdderessTemp, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.IntVar(&reportIntervalSecondTemp, "r", 10,
		"interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&pollIntervalSecondTemp, "p", 2,
		"interval in seconds for check metrics [default: 2 second]")
	flag.IntVar(&rateLimitTemp, "l", 1,
		"limit count reqyst in time [default: 1]")
	flag.StringVar(&signKeyTemp, "k", "",
		"key for sign [default: nil]")
	flag.StringVar(&pathToKey, "crypto-key", "",
		"path to file with key for [default: nil]")
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
	if reportIntervalSecondTemp != 0 {
		config.ReportIntervalSecond = reportIntervalSecondTemp
	}
	if pollIntervalSecondTemp != 0 {
		config.PollIntervalSecond = pollIntervalSecondTemp
	}
	if rateLimitTemp != 0 {
		config.RateLimit = rateLimitTemp
	}

	var err error
	if signKeyTemp != "" {
		config.SignKey, err = common.GetKeyFromString(signKeyTemp)
	}
	if err != nil {
		return fmt.Errorf("try get sign key: %w", err)
	}

	if envRateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		value, err := strconv.Atoi(envRateLimit)
		if err != nil {
			return fmt.Errorf("try atoi . value; err: %w", err)
		} else {
			config.RateLimit = value
		}
	}
	if envReportIntervalSecond, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		value, err := strconv.Atoi(envReportIntervalSecond)
		if err != nil {
			return fmt.Errorf("try atoi REPORT_INTERVAL value; err: %w", err)
		} else {
			config.ReportIntervalSecond = value
		}
	}
	if envPollIntervalSecond, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		value, err := strconv.Atoi(envPollIntervalSecond)
		if err != nil {
			return fmt.Errorf("try atoi POLL_INTERVAL value; err: %w", err)
		} else {
			config.PollIntervalSecond = value
		}
	}
	if envPathCryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		pathToKey = envPathCryptoKey
	}
	if pathToKey != "" {
		config.CryptoKeyOpen = common.GetDataFromFile(pathToKey)
	}
	return nil
}

func parseJson(data []byte) *Config {
	var conf configJSON
	err := json.Unmarshal(data, &conf)
	if err != nil {
		log.Printf("Unmarshal err. err:%s", err)
		return nil
	}
	var config Config
	config.CryptoKeyOpen = common.GetDataFromFile(conf.CryptoKeyOpen)
	config.ServerAdderess = conf.ServerAdderess
	config.SignKey, _ = common.GetKeyFromString(conf.SignKey)
	config.RateLimit = conf.RateLimit
	duration, err := time.ParseDuration(conf.PollIntervalSecond)
	if err != nil {
		log.Printf("ParseDuration 1 err. err:%s", err)
		return nil
	}
	config.PollIntervalSecond = int(duration.Seconds())
	duration, err = time.ParseDuration(conf.ReportIntervalSecond)
	if err != nil {
		log.Printf("ParseDuration 2 err. err:%s", err)
		return nil
	}
	config.ReportIntervalSecond = int(duration.Seconds())
	return &config
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
