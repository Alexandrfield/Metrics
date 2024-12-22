package server

import (
	"flag"
	"os"
)

func parseFlags(config *Config) {
	flag.StringVar(&config.ServerAdderess, "a", "localhost:8080",
		"address and port to run server [default:localhost:8080]")
	flag.Parse()

	if envServerAdderess, ok := os.LookupEnv("ADDRESS"); ok {
		config.ServerAdderess = envServerAdderess
	}
}

func GetServerConfig() Config {
	var config Config
	parseFlags(&config)
	return config
}
