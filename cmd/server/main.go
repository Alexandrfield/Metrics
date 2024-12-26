package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	"github.com/Alexandrfield/Metrics/internal/server"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Rcovert. Panic occurred: %w\n", err)
			debug.PrintStack()
		}
	}()
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can not initializate zap logger. err:%w", err)
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()

	config := server.GetServerConfig()
	stor := storage.CreateMemStorage()
	metricRep := server.CreateMetricRepository(stor, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	router := chi.NewRouter()
	router.Get(`/value/*`, server.WithLogging(logger, servHandler.GetValue))
	router.Get(`/value/`, server.WithLogging(logger, servHandler.GetJSONValue))
	router.Get(`/`, server.WithLogging(logger, servHandler.GetAllData))

	router.Post(`/update/*`, server.WithLogging(logger, servHandler.UpdateValue))
	router.Post(`/update/`, server.WithLogging(logger, servHandler.UpdateJSONValue))
	// router.Post(`/update/`, server.WithLogging(logger, servHandler.DefaultAnswer))

	logger.Info("Server stated")
	err = http.ListenAndServe(config.ServerAdderess, router)
	if err != nil {
		logger.Fatal("Unexpected error. err:%s", err)
	}
	logger.Info("Server stoped")
}
