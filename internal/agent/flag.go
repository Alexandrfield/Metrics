package agent

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func parseFlags(conf *Config) error {
	flag.StringVar(&conf.ServerAdderess, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.IntVar(&conf.ReportIntervalSecond, "r", 10,
		"interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&conf.PollIntervalSecond, "p", 2,
		"interval in seconds for check metrics [default: 2 second]")
	flag.Parse()

	if envServerAdderess := os.Getenv("ADDRESS"); envServerAdderess != "" {
		conf.ServerAdderess = envServerAdderess
	}

	if envReportIntervalSecond, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		value, err := strconv.Atoi(envReportIntervalSecond)
		if err != nil {
			return fmt.Errorf("Try atoi REPORT_INTERVAL value; err: %w\n", err)
		} else {
			conf.ReportIntervalSecond = value
		}
	}
	if envPollIntervalSecond, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		value, err := strconv.Atoi(envPollIntervalSecond)
		if err != nil {
			return fmt.Errorf("Try atoi POLL_INTERVAL value; err: %w\n", err)
		} else {
			conf.PollIntervalSecond = value
		}
	}
	return nil
}

func GetAgentConfig() (Config, error) {
	var config Config
	err := parseFlags(&config)
	if err != nil {
		return config, fmt.Errorf("GetAgentConfig err:%w", err)
	}
	return config, nil
}
