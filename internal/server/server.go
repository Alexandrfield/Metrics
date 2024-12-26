package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

//var errNotImplementedIssue = errors.New("not supported")

type MetricRepository struct {
	Logger       common.Loger
	LocalStorage *storage.MemStorage
}

func (rep *MetricRepository) SetCounterValue(metricName string, metricValue storage.TypeCounter) error {
	rep.Logger.Debugf("metricName:%s; metricValue:%s", metricName, metricValue)
	if rep.LocalStorage == nil {
		rep.Logger.Debugf("metricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	err := rep.LocalStorage.AddCounter(metricName, metricValue)
	if err != nil {
		return fmt.Errorf("problem SetCounterValue. err:%w", err)
	}
	return nil
}

func (rep *MetricRepository) SetGaugeValue(metricName string, metricValue storage.TypeGauge) error {
	rep.Logger.Debugf("metricName:%s; metricValue:%s\n", metricName, metricValue)
	if rep.LocalStorage == nil {
		rep.Logger.Debugf("metricRepository has not been initialize! Create default MemStorage\n")
		rep.LocalStorage = storage.CreateMemStorage()
	}
	err := rep.LocalStorage.AddGauge(metricName, metricValue)
	if err != nil {
		return fmt.Errorf("problem SetGaugeValue. err:%w", err)
	}
	return nil
}

func (rep *MetricRepository) GetCounterValue(metricName string) (storage.TypeCounter, error) {
	if rep.LocalStorage == nil {
		return storage.TypeCounter(0), errors.New("metricRepository has not been initialize")
	}
	val, err := rep.LocalStorage.GetCounter(metricName)
	if err != nil {
		return val, fmt.Errorf("problem GetCounterValue. err:%w", err)
	}
	return val, nil
}
func (rep *MetricRepository) GetGaugeValue(metricName string) (storage.TypeGauge, error) {
	if rep.LocalStorage == nil {
		return storage.TypeGauge(0), errors.New("metricRepository has not been initialize")
	}
	val, err := rep.LocalStorage.GetGauge(metricName)
	if err != nil {
		return val, fmt.Errorf("problem GetGaugeValue. err:%w", err)
	}
	return val, nil
}

func (rep *MetricRepository) GetAllValue() ([]string, error) {
	res := make([]string, 0)
	if rep.LocalStorage == nil {
		return res, errors.New("metricRepository has not been initialize")
	}
	allGaugeKeys, allCounterKeys := rep.LocalStorage.GetAllMetricName()
	for _, val := range allGaugeKeys {
		t, _ := rep.LocalStorage.GetGauge(val)
		res = append(res, fmt.Sprintf("name:%s; value:%v;\n", val, t))
	}
	for _, val := range allCounterKeys {
		t, _ := rep.LocalStorage.GetGauge(val)
		res = append(res, fmt.Sprintf("name:%s; value:%v;\n", val, t))
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
		logger.Infof("uri:%s; method:%s; status:%s; size:%s; duration:%s;",
			uri, method, responseData.status, responseData.size, duration)
	}
	return http.HandlerFunc(logFn)
}
