package requesthandler

import (
	"fmt"

	"github.com/Alexandrfield/Metrics/internal/storage"
)

var globalMemStorage *storage.MemStorage = nil

func HandleRequest(url []string) {
	fmt.Printf("url:%v\n", url)
	if globalMemStorage == nil {
		globalMemStorage = storage.CreateMemStorage()
	}
	switch url[2] {
	case "counter":
		globalMemStorage.AddCounter(url[3], url[4])
	case "gauge":
		globalMemStorage.AddGauge(url[3], url[4])
	}
}
