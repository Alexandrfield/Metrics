package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Alexandrfield/Metrics/internal/agent"
)

func parseFlags(conf *agent.Config) {
	flag.StringVar(&conf.ServerAdderess, "a", "localhost:8080", "address and port to run server [default:localhost:8080]")
	flag.IntVar(&conf.ReportIntervalSecond, "r", 10, "interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&conf.PollIntervalSecond, "p", 2, "interval in seconds for check metrics [default: 2 second]")
	flag.Parse()

	if envServerAdderess := os.Getenv("ADDRESS"); envServerAdderess != "" {
		conf.ServerAdderess = envServerAdderess
	}

	if envReportIntervalSecond, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		value, err := strconv.Atoi(envReportIntervalSecond)
		if err != nil {
			fmt.Printf("error atoi REPORT_INTERVAL; err: %s\n", err)
		} else {
			conf.ReportIntervalSecond = value
		}
	}
	if envPollIntervalSecond, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		value, err := strconv.Atoi(envPollIntervalSecond)
		if err != nil {
			fmt.Printf("error atoi POLL_INTERVAL; err: %s\n", err)
		} else {
			conf.PollIntervalSecond = value
		}
	}

}

func GetAgentConfig() agent.Config {
	var conf agent.Config
	parseFlags(&conf)
	return conf
}
