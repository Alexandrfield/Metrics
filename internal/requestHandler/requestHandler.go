package requesthandler

import (
	"fmt"

	"github.com/Alexandrfield/Metrics/internal/storage"
)

var globalMemStorage storage.MemStorageI = nil

func HandleRequest(url []string) bool {
	status := false
	fmt.Printf("url:%v\n", url)
	if globalMemStorage == nil {
		globalMemStorage = storage.CreateMemStorage()
	}
	switch url[2] {
	case "counter":
		status = globalMemStorage.AddCounter(url[3], url[4])
	case "gauge":
		status = globalMemStorage.AddGauge(url[3], url[4])
	}
	return status
}
