package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexandrfield/Metrics/internal/agent"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can not initializate zap logger. err:%w", err)
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()

	agentConfig := agent.GetAgentConfig()
	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}

	done := make(chan struct{})
	go agent.MetricsWatcher(agentConfig, &client, logger, done)
	logger.Info("Agent stated")
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
	close(done)
	time.Sleep(1 * time.Second)
	logger.Info("Agent stoped")
}
