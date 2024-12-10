package storage

import (
	"errors"
	"fmt"
	"strconv"
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

func (st *MemStorage) AddGauge(name string, raw string) error {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("error parse Gauge type. Value:%s; error parse :%w;", raw, err)
	}
	st.gaugeData[name] = TypeGauge(value)
	return nil
}
func (st *MemStorage) GetGauge(name string) (string, error) {
	val, ok := st.gaugeData[name]
	if !ok {
		return "", fmt.Errorf("can't find Gauge metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return strconv.FormatFloat(float64(val), 'f', -1, 64), nil
}
func (st *MemStorage) AddCounter(name string, raw string) error {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("error parse Counter type. Value:%s; err parse:%w;", raw, err)
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
		return "", fmt.Errorf("can't find Counter metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return strconv.Itoa(int(val)), nil
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
