package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

func main() {
	config := GetServerConfig()
	storage := storage.CreateMemStorage()
	servHandler := handler.CreateHandlerRepository(storage)
	router := chi.NewRouter()
	router.Get(`/value/*`, servHandler.getValue)
	router.Get(`/`, servHandler.getAllData)

	router.Post(`/update/*`, servHandler.updateValue)
	router.Post(`/update/`, servHandler.defaultAnswer)

	err := http.ListenAndServe(config.ServerAdderess, router)
	if err != nil {
		log.Fatal("Unexpected error. err:%w", err)
	}
}
