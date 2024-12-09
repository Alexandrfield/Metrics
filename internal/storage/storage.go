package storage

import (
	"fmt"
	"strconv"
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

func (st *MemStorage) AddGauge(name string, d string) bool {
	value, err := strconv.ParseFloat(d, 64)
	if err != nil {
		fmt.Printf("error parse gauge; err: %s\n", err)
		return false
	}
	st.gaugeData[name] = TypeGauge(value)
	return true
}
func (st *MemStorage) GetGauge(name string) (string, bool) {
	val, ok := st.gaugeData[name]
	res := ""
	if ok {
		res = strconv.FormatFloat(float64(val), 'f', -1, 64)
	}
	return res, ok
}
func (st *MemStorage) AddCounter(name string, d string) bool {
	value, err := strconv.Atoi(d)
	if err != nil {
		fmt.Printf("error parse counter; err: %s\n", err)
		return false
	}
	val, ok := st.counterData[name]
	if !ok {
		val = 0
	}
	st.counterData[name] = val + TypeCounter(value)
	return true
}
func (st *MemStorage) GetCounter(name string) (string, bool) {
	val, ok := st.counterData[name]
	res := ""
	if ok {
		res = fmt.Sprintf("%d", val)
	}
	return res, ok
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
