package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexandrfield/Metrics/internal/agent"
)

func main() {
	agentConfig := agent.GetAgentConfig()
	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}

	done := make(chan struct{})
	go agent.MetricsWatcher(agentConfig, &client, done)
	log.Println("Agent stated")
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
	close(done)
	time.Sleep(1 * time.Second)
}
