package storage

import (
	"fmt"
	"strconv"
)

type TypeGauge float64
type TypeCounter int64

type MemStorageI interface {
	AddGauge(name string, d string) bool
	GetGauge(name string) (TypeGauge, bool)
	AddCounter(name string, d string) bool
	GetCounter(name string) (TypeCounter, bool)
}
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
func (st *MemStorage) GetGauge(name string) (TypeGauge, bool) {
	val, ok := st.gaugeData[name]
	return val, ok
}
func (st *MemStorage) AddCounter(name string, d string) bool {
	value, err := strconv.Atoi(d)
	if err != nil {
		fmt.Printf("error parse counter; err: %s\n", err)
		return false
	}
	val, ok := st.counterData[name]
	if ok {
		val = 0
	}
	st.counterData[name] = val + TypeCounter(value)
	return true
}
func (st *MemStorage) GetCounter(name string) (TypeCounter, bool) {
	val, ok := st.counterData[name]
	return val, ok
}
