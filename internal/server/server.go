package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"log"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

var errNotImplementedIssue = errors.New("not supported")

type MetricRepository struct {
	LocalStorage *storage.MemStorage
}

func (rep *MetricRepository) SetCounterValue(metricName string, metricValue storage.TypeCounter) error {
	log.Printf("metricName:%s; metricValue:%s\n", metricName, metricValue)
	if rep.LocalStorage == nil {
		log.Printf("metricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	return rep.LocalStorage.AddCounter(metricName, metricValue)
}

func (rep *MetricRepository) SetGaugeValue(metricName string, metricValue storage.TypeGauge) error {
	log.Printf("metricName:%s; metricValue:%s\n", metricName, metricValue)
	if rep.LocalStorage == nil {
		log.Printf("metricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	return rep.LocalStorage.AddGauge(metricName, metricValue)
}

func (rep *MetricRepository) GetCounterValue(metricName string) (storage.TypeCounter, error) {
	if rep.LocalStorage == nil {
		return storage.TypeCounter(0), errors.New("metricRepository has not been initialize")
	}
	return rep.LocalStorage.GetCounter(metricName)
}
func (rep *MetricRepository) GetGaugeValue(metricName string) (storage.TypeGauge, error) {
	if rep.LocalStorage == nil {
		return storage.TypeGauge(0), errors.New("metricRepository has not been initialize")
	}
	return rep.LocalStorage.GetGauge(metricName)
}

func (rep *MetricRepository) GetAllValue() ([]string, error) {
	res := make([]string, 0)
	if rep.LocalStorage == nil {
		return res, errors.New("metricRepository has not been initialize")
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

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(logger common.Loger, h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logger.Infof("uri:%s; method:%s; status:%s; size:%s; duration:%s;", uri, method, responseData.status, responseData.size, duration)

	}
	return http.HandlerFunc(logFn)
}
