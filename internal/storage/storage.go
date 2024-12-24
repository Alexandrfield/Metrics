package storage

import (
	"errors"
	"fmt"
)

var errMetricNotExistIssue = errors.New("metric with this name or type is does't exist")

type TypeGauge float64
type TypeCounter int64

type MemStorage struct {
	gaugeData   map[string]TypeGauge
	counterData map[string]TypeCounter
}

func CreateMemStorage() *MemStorage {
	return &MemStorage{gaugeData: make(map[string]TypeGauge), counterData: make(map[string]TypeCounter)}
}

func (st *MemStorage) AddGauge(name string, value TypeGauge) error {
	st.gaugeData[name] = value
	return nil
}
func (st *MemStorage) GetGauge(name string) (TypeGauge, error) {
	val, ok := st.gaugeData[name]
	if !ok {
		return TypeGauge(0), fmt.Errorf("can't find Gauge metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemStorage) AddCounter(name string, value TypeCounter) error {
	val, ok := st.counterData[name]
	if !ok {
		val = 0
	}
	st.counterData[name] = val + TypeCounter(value)
	return nil
}
func (st *MemStorage) GetCounter(name string) (TypeCounter, error) {
	val, ok := st.counterData[name]
	if !ok {
		return TypeCounter(0), fmt.Errorf("can't find Counter metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	for key := range st.gaugeData {
		allGaugeKeys = append(allGaugeKeys, key)
	}
	allCounterKeys := make([]string, 0)
	for key := range st.counterData {
		allCounterKeys = append(allCounterKeys, key)
	}
	return allGaugeKeys, allCounterKeys
}
