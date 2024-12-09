package requesthandler

import (
	"fmt"

	"github.com/Alexandrfield/Metrics/internal/storage"
)

type CommonMemStorage interface {
	AddGauge(name string, d string) bool
	GetGauge(name string) (string, bool)
	AddCounter(name string, d string) bool
	GetCounter(name string) (string, bool)
	GetAllMetricName() ([]string, []string)
}

type Repository struct {
	memStorage CommonMemStorage
}

func CreateHandlerRepository(stor CommonMemStorage) *Repository {
	return &Repository{memStorage: stor}
}

func (rep *Repository) HandleRequest(url []string) bool {
	status := false
	fmt.Printf("url:%v\n", url)
	if rep.memStorage == nil {
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
func HandleGetValue(url []string) (string, bool) {
	status := false
	res := ""
	if globalMemStorage == nil {
		return "", false
	}
	switch url[2] {
	case "counter":
		res, status = globalMemStorage.GetCounter(url[3])
	case "gauge":
		res, status = globalMemStorage.GetGauge(url[3])
	}
	return res, status
}
func HandleAllValue() []string {
	var res []string
	if globalMemStorage == nil {
		return res
	}
	allGaugeKeys, allCounterKeys := globalMemStorage.GetAllMetricName()
	for i := 0; i < len(allGaugeKeys); i++ {
		t, _ := globalMemStorage.GetGauge(allGaugeKeys[i])
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", allGaugeKeys[i], t))
	}
	for i := 0; i < len(allCounterKeys); i++ {
		t, _ := globalMemStorage.GetGauge(allCounterKeys[i])
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", allCounterKeys[i], t))
	}
	return res
}
