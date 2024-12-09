package server

import (
	"errors"
	"fmt"

	"log"

	"github.com/Alexandrfield/Metrics/internal/storage"
)

var ErrNotImplementedIssue = errors.New("not supported")

type MetricRepository struct {
	LocalStorage *storage.MemStorage
}

func (rep *MetricRepository) SetValue(metricType string, metricName string, metricValue string) error {
	var err error
	log.Printf("metricType:%s; metricValue:%s\n", metricType, metricValue)
	if rep.LocalStorage == nil {
		log.Printf("metricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	switch metricType {
	case "counter":
		err = rep.LocalStorage.AddCounter(metricName, metricValue)
	case "gauge":
		err = rep.LocalStorage.AddGauge(metricName, metricValue)
	default:
		err = fmt.Errorf("unknown type %s;err:%w", metricType, ErrNotImplementedIssue)
	}
	return err
}

func (rep *MetricRepository) GetValue(metricType string, metricName string) (string, error) {
	var err error
	res := ""
	if rep.LocalStorage == nil {
		return res, fmt.Errorf("metricRepository has not been initialize. err:%w", storage.ErrMetricNotExistIssue)
	}
	switch metricType {
	case "counter":
		log.Printf("counter -> metricName:%s\n", metricName)
		res, err = rep.LocalStorage.GetCounter(metricName)
	case "gauge":
		log.Printf("gauge -> metricName:%s\n", metricName)
		res, err = rep.LocalStorage.GetGauge(metricName)
	default:
		err = fmt.Errorf("unknown type %s; err:%w", metricType, ErrNotImplementedIssue)
	}
	return res, err
}

func (rep *MetricRepository) GetAllValue() ([]string, error) {
	res := make([]string, 0)
	if rep.LocalStorage == nil {
		return res, fmt.Errorf("metricRepository has not been initialize. err:%w", storage.ErrMetricNotExistIssue)
	}
	allGaugeKeys, allCounterKeys := rep.LocalStorage.GetAllMetricName()
	for _, val := range allGaugeKeys {
		t, _ := rep.LocalStorage.GetGauge(val)
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", val, t))
	}
	for _, val := range allCounterKeys {
		t, _ := rep.LocalStorage.GetGauge(val)
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", val, t))
	}
	return res, nil
}
