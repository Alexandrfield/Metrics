package server

import (
	"fmt"

	"github.com/Alexandrfield/Metrics/internal/customErrors"
	"github.com/Alexandrfield/Metrics/internal/storage"
	"gvisor.dev/gvisor/pkg/log"
)

type MetricRepository struct {
	LocalStorage *storage.MemStorage
}

func (rep *MetricRepository) SetValue(metricType string, metricName string, metricValue string) error {
	var err error
	log.Debugf("metricType:%s; metricValue:%s\n", metricType, metricValue)
	if rep.LocalStorage == nil {
		log.Infof("MetricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	switch metricType {
	case "counter":
		err = rep.LocalStorage.AddCounter(metricName, metricValue)
	case "gauge":
		err = rep.LocalStorage.AddGauge(metricName, metricValue)
	}
	return err
}

func (rep *MetricRepository) GetValue(metricType string, metricName string) (string, error) {
	var err error
	res := ""
	if rep.LocalStorage == nil {
		return res, fmt.Errorf("MetricRepository has not been initialize. err:%w", customErrors.ErrMetricNotExistIssue)
	}
	switch metricType {
	case "counter":
		res, err = rep.LocalStorage.GetCounter(metricName)
	case "gauge":
		res, err = rep.LocalStorage.GetGauge(metricName)
	}
	return res, err
}

func (rep *MetricRepository) GetAllValue() ([]string, error) {
	var res []string
	if rep.LocalStorage == nil {
		return res, fmt.Errorf("MetricRepository has not been initialize. err:%w", customErrors.ErrMetricNotExistIssue)
	}
	allGaugeKeys, allCounterKeys := rep.LocalStorage.GetAllMetricName()
	for i := 0; i < len(allGaugeKeys); i++ {
		t, _ := rep.LocalStorage.GetGauge(allGaugeKeys[i])
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", allGaugeKeys[i], t))
	}
	for i := 0; i < len(allCounterKeys); i++ {
		t, _ := rep.LocalStorage.GetGauge(allCounterKeys[i])
		res = append(res, fmt.Sprintf("name:%s; value:%s;\n", allCounterKeys[i], t))
	}
	return res, nil
}
