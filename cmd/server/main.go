package main

import (
	"log"
	"net/http"

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
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	config := server.GetServerConfig()
	stor := storage.CreateMemStorage()
	metricRep := server.MetricRepository{LocalStorage: stor, Logger: logger}
	servHandler := handler.CreateHandlerRepository(&metricRep)

	router := chi.NewRouter()
	router.Get(`/value/*`, server.WithLogging(logger, servHandler.GetValue))
	router.Get(`/`, server.WithLogging(logger, servHandler.GetAllData))

	router.Post(`/update/*`, server.WithLogging(logger, servHandler.UpdateValue))
	router.Post(`/update/`, server.WithLogging(logger, servHandler.DefaultAnswer))

	logger.Info("Server stated")
	err = http.ListenAndServe(config.ServerAdderess, router)
	if err != nil {
		logger.Fatal("Unexpected error. err:%s", err)
	}
	logger.Info("Server stoped")
}
