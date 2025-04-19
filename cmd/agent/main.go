package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexandrfield/Metrics/internal/agent"
	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can not initializate zap logger. err:%w", err)
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()

	agentConfig, err := agent.GetAgentConfig()
	if err != nil {
		logger.Fatalf("Cant inint agent. err:%w", err)
	}
	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}

	done := make(chan struct{})
	defer func() {
		logger.Info("Agent stoping...")
		close(done)
		time.Sleep(1 * time.Second)
		logger.Info("Agent stoped")
	}()
	go agent.MetricsWatcher(agentConfig, &client, logger, done)
	logger.Info("Agent started")
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
}
