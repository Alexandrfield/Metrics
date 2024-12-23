package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	"github.com/Alexandrfield/Metrics/internal/server"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

func main() {
	config := server.GetServerConfig()
	stor := storage.CreateMemStorage()
	metricRep := server.MetricRepository{LocalStorage: stor}
	servHandler := handler.CreateHandlerRepository(&metricRep)

	router := chi.NewRouter()
	router.Get(`/value/*`, servHandler.GetValue)
	router.Get(`/`, servHandler.GetAllData)

	router.Post(`/update/*`, servHandler.UpdateValue)
	router.Post(`/update/`, servHandler.DefaultAnswer)

	log.Println("Server stated")
	err := http.ListenAndServe(config.ServerAdderess, router)
	if err != nil {
		log.Fatal("Unexpected error. err:%w", err)
	}
}
