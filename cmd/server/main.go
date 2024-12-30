package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

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
			fmt.Printf("Rcovert. Panic occurred:\n")
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	config := server.GetServerConfig()
	done := make(chan struct{})
	logger.Debugf("config file: %v", config)
	storageConfig := storage.Config{FileStoregePath: config.FileStoregePath,
		StoreIntervalSecond: config.StoreIntervalSecond, Restore: config.Restore}
	stor := storage.CreateMemStorage(storageConfig, logger, done)
	if stor == nil {
		logger.Fatal("Can not create MemStorage. err:%s", err)
	}
	metricRep := server.CreateMetricRepository(stor, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	router := chi.NewRouter()
	router.Get(`/value/*`, server.WithLogging(logger, servHandler.GetValue))
	router.Post(`/value/`, server.WithLogging(logger, servHandler.GetJSONValue))
	router.Get(`/`, server.WithLogging(logger, servHandler.GetAllData))

	router.Post(`/update/*`, server.WithLogging(logger, servHandler.UpdateValue))
	router.Post(`/update/`, server.WithLogging(logger, servHandler.UpdateJSONValue))
	// router.Post(`/update/`, server.WithLogging(logger, servHandler.DefaultAnswer))

	logger.Info("Server started")
	defer func() {
		logger.Info("Server stoping ... ")
		close(done)
		time.Sleep(1 * time.Second)
		logger.Info("Server stoped")
	}()
	err = http.ListenAndServe(config.ServerAdderess, router)
	if err != nil {
		logger.Fatal("Unexpected error. err:%s", err)
	}
}
