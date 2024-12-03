package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var globalPollIntervalSecond int
var globalReportIntervalSecond int
var globalServerAdderess string

func parseFlags() {
	flag.StringVar(&globalServerAdderess, "a", "localhost:8080", "address and port to run server [default:localhost:8080]")
	flag.IntVar(&globalReportIntervalSecond, "r", 10, "interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&globalPollIntervalSecond, "p", 2, "interval in seconds for check metrics [default: 2 second]")
	flag.Parse()

	if envServerAdderess := os.Getenv("ADDRESS"); envServerAdderess != "" {
		globalServerAdderess = envServerAdderess
	}

	if envReportIntervalSecond := os.Getenv("REPORT_INTERVAL"); envReportIntervalSecond != "" {
		value, err := strconv.Atoi(envReportIntervalSecond)
		if err != nil {
			fmt.Printf("error atoi REPORT_INTERVAL; err: %s\n", err)
		} else {
			globalReportIntervalSecond = value
		}
	}
	if envPollIntervalSecond := os.Getenv("POLL_INTERVAL"); envPollIntervalSecond != "" {
		value, err := strconv.Atoi(envPollIntervalSecond)
		if err != nil {
			fmt.Printf("error atoi POLL_INTERVAL; err: %s\n", err)
		} else {
			globalPollIntervalSecond = value
		}
	}

}
