package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	_ "net/http/pprof"

	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	"github.com/Alexandrfield/Metrics/internal/server"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can not initializate zap logger. err:%w", err)
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Rcovert. Panic occurred. err:%w", err)
			debug.PrintStack()
		}
	}()

	go func() {
		logger.Infof("start ", http.ListenAndServe("localhost:6060", nil))
	}()

	config, err := server.GetServerConfig()
	if err != nil {
		logger.Fatalf("Cant init server. err:%w", err)
	}
	done := make(chan struct{})
	defer func() {
		close(done)
		logger.Info("Server stoping ... ")
		time.Sleep(1 * time.Second)
		logger.Info("Server stoped")
	}()
	logger.Debugf("config file ServerAdderess: %s; FileStoregePath:%s; database:",
		config.ServerAdderess, config.FileStoregePath)
	storageConfig := storage.Config{FileStoregePath: config.FileStoregePath,
		StoreIntervalSecond: config.StoreIntervalSecond, Restore: config.Restore, DatabaseDsn: config.DatabaseDsn}
	stor := storage.CreateMemStorage(storageConfig, logger, done)
	if stor == nil {
		logger.Fatal("Can not create MemStorage. err:%s", err)
	}
	metricRep := server.CreateMetricRepository(stor, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	router := chi.NewRouter()
	router.Get(`/value/*`, server.WithLogging(logger, &config, servHandler.GetValue))
	router.Post(`/value/`, server.WithLogging(logger, &config, servHandler.GetJSONValue))
	router.Get(`/`, server.WithLogging(logger, &config, servHandler.GetAllData))

	router.Get(`/ping`, server.WithLogging(logger, &config, servHandler.Ping))

	router.Post(`/update/*`, server.WithLogging(logger, &config, servHandler.UpdateValue))
	router.Post(`/update/`, server.WithLogging(logger, &config, servHandler.UpdateJSONValue))
	router.Post(`/updates/`, server.WithLogging(logger, &config, servHandler.UpdatesMetrics))

	logger.Info("Server started")
	go func() {
		err = http.ListenAndServe(config.ServerAdderess, router)
		if err != nil {
			logger.Errorf("Unexpected error. err:%s", err)
		}
	}()
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
}
