package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

func updateGaugeMetrics(metrics map[string]storage.TypeGauge) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	metrics["Alloc"] = storage.TypeGauge(rtm.Alloc)
	metrics["BuckHashSys"] = storage.TypeGauge(rtm.BuckHashSys)
	metrics["Frees"] = storage.TypeGauge(rtm.Frees)
	metrics["GCCPUFraction"] = storage.TypeGauge(rtm.GCCPUFraction)
	metrics["GCSys"] = storage.TypeGauge(rtm.GCSys)
	metrics["HeapAlloc"] = storage.TypeGauge(rtm.HeapAlloc)
	metrics["HeapIdle"] = storage.TypeGauge(rtm.HeapIdle)
	metrics["HeapInuse"] = storage.TypeGauge(rtm.HeapInuse)
	metrics["HeapObjects"] = storage.TypeGauge(rtm.HeapObjects)
	metrics["HeapReleased"] = storage.TypeGauge(rtm.HeapReleased)
	metrics["HeapSys"] = storage.TypeGauge(rtm.HeapSys)
	metrics["LastGC"] = storage.TypeGauge(rtm.LastGC)
	metrics["Lookups"] = storage.TypeGauge(rtm.Lookups)
	metrics["MCacheInuse"] = storage.TypeGauge(rtm.MCacheInuse)
	metrics["MCacheSys"] = storage.TypeGauge(rtm.MCacheSys)
	metrics["MSpanInuse"] = storage.TypeGauge(rtm.MSpanInuse)
	metrics["MSpanSys"] = storage.TypeGauge(rtm.MSpanSys)
	metrics["Mallocs"] = storage.TypeGauge(rtm.Mallocs)
	metrics["NextGC"] = storage.TypeGauge(rtm.NextGC)
	metrics["NumForcedGC"] = storage.TypeGauge(rtm.NumForcedGC)
	metrics["NumGC"] = storage.TypeGauge(rtm.NumGC)
	metrics["OtherSys"] = storage.TypeGauge(rtm.OtherSys)
	metrics["PauseTotalNs"] = storage.TypeGauge(rtm.PauseTotalNs)
	metrics["StackInuse"] = storage.TypeGauge(rtm.StackInuse)
	metrics["StackSys"] = storage.TypeGauge(rtm.StackSys)
	metrics["Sys"] = storage.TypeGauge(rtm.Sys)
	metrics["TotalAlloc"] = storage.TypeGauge(rtm.TotalAlloc)
	metrics["RandomValue"] = storage.TypeGauge(rand.Float64())
}
func updateCounterMetrics(metrics map[string]storage.TypeCounter) {
	metrics["PollCount"]++
}
func prepareReportGaugeMetrics(metricsGauge map[string]storage.TypeGauge) []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metricsGauge {
		temp := float64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "gauge", Value: &temp})
	}
	return dataMetricForReport
}

func prepareReportCounterMetrics(metricsCounter map[string]storage.TypeCounter) []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metricsCounter {
		temp := int64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "counter", Delta: &temp})
	}
	return dataMetricForReport
}

func reportMetrics(client *http.Client, serverAdderess string, dataMetricForReport []common.Metrics,
	logger common.Loger) {
	for _, metric := range dataMetricForReport {
		_, err := reportMetric(client, serverAdderess, metric, logger)
		if err != nil {
			logger.Warnf("error report metric. err%s\n ", err)
		}
	}
}
func reportMetric(client *http.Client, serverAdderess string, metric common.Metrics,
	logger common.Loger) (int, error) {
	objMetrics, err := json.Marshal(metric)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("problem with marshal JSON file. err:%w", err)
	}
	url := fmt.Sprintf("http://%s/update/", serverAdderess)

	req, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(objMetrics),
	)
	if err != nil {
		logger.Warnf("http.NewRequest. err: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Warnf("http.NewRequest.Do err: %s\n", err)
		return 0, fmt.Errorf("http.NewRequest.Do err:%w", err)
	}
	status := resp.StatusCode
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Warnf("resp.Body.Close() err: %s\n", err)
		}
	}()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return status, fmt.Errorf("error reading body. err:%w", err)
	}
	return status, nil
}

func reportCounterMetrics(client *http.Client, serverAdderess string, dataMetricForReport []common.Metrics,
	metricsCounter map[string]storage.TypeCounter, logger common.Loger) {
	for _, metric := range dataMetricForReport {
		statusCode, err := reportMetric(client, serverAdderess, metric, logger)
		if err != nil {
			logger.Warnf("error report metric for counter. err%s\n ", err)
			continue
		}
		if statusCode == http.StatusOK {
			metricsCounter[metric.ID] -= storage.TypeCounter(*metric.Delta)
		}
	}
}
func MetricsWatcher(config Config, client *http.Client, logger common.Loger, done chan struct{}) {
	tickerPoolInterval := time.NewTicker(time.Duration(config.PollIntervalSecond) * time.Second)
	tickerReportInterval := time.NewTicker(time.Duration(config.ReportIntervalSecond) * time.Second)
	metricsGauge := make(map[string]storage.TypeGauge)
	metricsCounter := make(map[string]storage.TypeCounter)
	metricsCounter["PollCount"] = 0
	for {
		select {
		case <-done:
			return
		case <-tickerPoolInterval.C:
			updateGaugeMetrics(metricsGauge)
			updateCounterMetrics(metricsCounter)
		case <-tickerReportInterval.C:
			metricsGaugeReport := prepareReportGaugeMetrics(metricsGauge)
			metricsCounterReport := prepareReportCounterMetrics(metricsCounter)
			reportMetrics(client, config.ServerAdderess, metricsGaugeReport, logger)
			reportCounterMetrics(client, config.ServerAdderess, metricsCounterReport, metricsCounter, logger)
		}
	}
}
