package server

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

// MetricRepository that object for processing metrics (Set, Get...).
type MetricRepository struct {
	Logger       common.Loger
	LocalStorage storage.BasicStorage
}

func CreateMetricRepository(localStorage storage.BasicStorage, logger common.Loger) MetricRepository {
	return MetricRepository{Logger: logger, LocalStorage: localStorage}
}

func (rep *MetricRepository) SetCounterValue(metricName string, metricValue common.TypeCounter) error {
	rep.Logger.Debugf("metricName:%s; metricValue:%s", metricName, metricValue)
	if rep.LocalStorage == nil {
		return errors.New("localStorage for repository not init")
	}
	err := rep.LocalStorage.AddCounter(metricName, metricValue)
	if err != nil {
		return fmt.Errorf("problem SetCounterValue. err:%w", err)
	}
	return nil
}

func (rep *MetricRepository) SetGaugeValue(metricName string, metricValue common.TypeGauge) error {
	rep.Logger.Debugf("metricName:%s; metricValue:%s\n", metricName, metricValue)
	if rep.LocalStorage == nil {
		return errors.New("localStorage for repository not init")
	}
	err := rep.LocalStorage.AddGauge(metricName, metricValue)
	if err != nil {
		return fmt.Errorf("problem SetGaugeValue. err:%w", err)
	}
	return nil
}

func (rep *MetricRepository) GetCounterValue(metricName string) (common.TypeCounter, error) {
	if rep.LocalStorage == nil {
		return common.TypeCounter(0), errors.New("metricRepository has not been initialize")
	}
	val, err := rep.LocalStorage.GetCounter(metricName)
	if err != nil {
		return val, fmt.Errorf("problem GetCounterValue. err:%w", err)
	}
	return val, nil
}
func (rep *MetricRepository) GetGaugeValue(metricName string) (common.TypeGauge, error) {
	if rep.LocalStorage == nil {
		return common.TypeGauge(0), errors.New("metricRepository has not been initialize")
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
func (rep *MetricRepository) PingDatabase() bool {
	rep.Logger.Debugf("PingDatabase")
	return rep.LocalStorage.PingDatabase()
}

func (rep *MetricRepository) AddMetrics(metrics []common.Metrics) error {
	rep.Logger.Debugf("AddMetrics len(metrics):%d", len(metrics))
	err := rep.LocalStorage.AddMetrics(metrics)
	if err != nil {
		return fmt.Errorf("problem wit add metrics. err:%w", err)
	}
	return nil
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
	if err != nil {
		err = fmt.Errorf("loggingResponseWriter err:%w", err)
	}
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(logger common.Loger, config *Config, h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		var lw loggingResponseWriter
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			logger.Debugf("try use gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				_, _ = io.WriteString(w, err.Error())
				logger.Debugf("gzip.NewWriterLevel error:%w", err)
			}
			w.Header().Set("Content-Encoding", "gzip")
			defer func() {
				_ = gz.Close()
			}()
			lw = loggingResponseWriter{
				ResponseWriter: gzipWriter{ResponseWriter: w, Writer: gz}, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
		} else {
			logger.Debugf("not use gzip")
			lw = loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
		}
		data := make([]byte, 10000)
		n, _ := r.Body.Read(data)
		data = data[:n]
		msgSign := r.Header.Get("HashSHA256")
		if msgSign != "" && len(data) > 0 {
			sig, _ := b64.StdEncoding.DecodeString(msgSign)
			if !common.CheckHash(data, sig, config.SignKey) {
				lw.WriteHeader(http.StatusBadRequest)
			}
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logger.Infof("uri:%s; method:%s; status:%d; size:%d; duration:%s;",
			uri, method, responseData.status, responseData.size, duration)
	}
	return http.HandlerFunc(logFn)
}
