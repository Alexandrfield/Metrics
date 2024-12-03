package main

import "flag"

var globalPollIntervalSecond int
var globalReportIntervalSecond int
var globalServerAdderess string

func parseFlags() {
	flag.StringVar(&globalServerAdderess, "a", "localhost:8080", "address and port to run server [default:localhost:8080]")
	flag.IntVar(&globalReportIntervalSecond, "r", 10, "interval in seconds  for sending report to server [default: 10 second]")
	flag.IntVar(&globalPollIntervalSecond, "p", 2, "interval in seconds for check metrics [default: 2 second]")
	flag.Parse()
}
