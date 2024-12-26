package common

import (
	"fmt"
	"strconv"
)

type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

func (met *Metrics) SaveMetric(mtype string, name string, rawValue string) error {
	met.ID = name
	met.MType = mtype
	switch met.MType {
	case "counter":
		value, err := strconv.Atoi(rawValue)
		if err != nil {
			return fmt.Errorf("error parse Counter type. Value:%s; err parse:%w;", rawValue, err)
		}
		temp := int64(value)
		met.Delta = &temp
	case "gauge":
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return fmt.Errorf("error parse Gauge type. Value:%s; error parse :%w;", rawValue, err)
		}
		met.Value = &value
	}
	return nil
}

func (met *Metrics) GetValueMetric() string {
	var res string
	switch met.MType {
	case "counter":
		res = strconv.Itoa(int(*met.Delta))
	case "gauge":
		res = strconv.FormatFloat(float64(*met.Value), 'f', -1, 64)
	}
	return res
}
