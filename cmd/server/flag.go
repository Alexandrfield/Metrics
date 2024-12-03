package main

import (
	"flag"
	"os"
)

var globalServerAdderess string

func parseFlags() {
	flag.StringVar(&globalServerAdderess, "a", "localhost:8080", "address and port to run server [default:localhost:8080]")
	flag.Parse()

	if envServerAdderess := os.Getenv("ADDRESS"); envServerAdderess != "" {
		globalServerAdderess = envServerAdderess
	}

}
