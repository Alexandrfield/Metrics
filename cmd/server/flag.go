package main

import (
	"flag"
	"os"

	"github.com/Alexandrfield/Metrics/internal/server"
)

func parseFlags(config *server.Config) {
	flag.StringVar(&config.ServerAdderess, "a", "localhost:8080", "address and port to run server [default:localhost:8080]")
	flag.Parse()

	if envServerAdderess, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAdderess = envServerAdderess
	}

}

func GetServerConfig() server.Config {
	var config server.Config
	parseFlags(&config)
	return config
}
