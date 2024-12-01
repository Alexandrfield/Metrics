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

func (st *MemStorage) AddGauge(name string, d string) {
	value, err := strconv.ParseFloat(d, 64)
	if err != nil {
		fmt.Printf("error parse gauge; err: %w\n", err)
		return
	}
	st.gaugeData[name] = TypeGauge(value)
}
func (st *MemStorage) AddCounter(name string, d string) {
	value, err := strconv.Atoi(d)
	if err != nil {
		fmt.Printf("error parse counter; err: %w\n", err)
		return
	}
	val, ok := st.counterData[name]
	if ok {
		val = 0
	}
	st.counterData[name] = val + TypeCounter(value)
}
