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
		err = fmt.Errorf("unknown type %s;err:%w", metricType, errNotImplementedIssue)
	}
	return err
}

func (rep *MetricRepository) GetValue(metricType string, metricName string) (string, error) {
	var err error
	res := ""
	if rep.LocalStorage == nil {
		return res, errors.New("metricRepository has not been initialize")
	}
	switch metricType {
	case "counter":
		log.Printf("counter -> metricName:%s\n", metricName)
		res, err = rep.LocalStorage.GetCounter(metricName)
	case "gauge":
		log.Printf("gauge -> metricName:%s\n", metricName)
		res, err = rep.LocalStorage.GetGauge(metricName)
	default:
		err = fmt.Errorf("unknown type %s; err:%w", metricType, errNotImplementedIssue)
	}
	return res, err
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

func WithLogging(logger common.Loger, h http.Handler) http.HandlerFunc {
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
