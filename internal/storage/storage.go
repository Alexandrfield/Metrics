package storage

import (
	"fmt"
	"strconv"

	"github.com/Alexandrfield/Metrics/internal/customErrors"
)

type TypeGauge float64
type TypeCounter int64

type MemStorage struct {
	gaugeData   map[string]TypeGauge
	counterData map[string]TypeCounter
}

func CreateMemStorage() *MemStorage {
	return &MemStorage{gaugeData: make(map[string]TypeGauge), counterData: make(map[string]TypeCounter)}
}

func (st *MemStorage) AddGauge(name string, raw string) error {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("Error parse Gauge type. Value:%s; err:%w", raw, customErrors.ErrCantParseDataIssue)
	}
	st.gaugeData[name] = TypeGauge(value)
	return nil
}
func (st *MemStorage) GetGauge(name string) (string, error) {
	val, ok := st.gaugeData[name]
	if !ok {
		return "", fmt.Errorf("Can't find Gauge metric with name:%s;err:%w", name, customErrors.ErrMetricNotExistIssue)
	}
	return strconv.FormatFloat(float64(val), 'f', -1, 64), nil
}
func (st *MemStorage) AddCounter(name string, raw string) error {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("Error parse Counter type. Value:%s; err:%w", raw, customErrors.ErrCantParseDataIssue)
	}
	val, ok := st.counterData[name]
	if !ok {
		val = 0
	}
	st.counterData[name] = val + TypeCounter(value)
	return nil
}
func (st *MemStorage) GetCounter(name string) (string, error) {
	val, ok := st.counterData[name]
	if !ok {
		return "", fmt.Errorf("Can't find Counter metric with name:%s;err:%w", name, customErrors.ErrMetricNotExistIssue)
	}
	return strconv.Itoa(int(val)), nil
}
func (st *MemStorage) GetAllMetricName() ([]string, []string) {
	var allGaugeKeys []string
	for key := range st.gaugeData {
		allGaugeKeys = append(allGaugeKeys, key)
	}
	var allCounterKeys []string
	for key := range st.counterData {
		allCounterKeys = append(allCounterKeys, key)
	}
	return allGaugeKeys, allCounterKeys
}
