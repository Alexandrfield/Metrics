package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"syscall"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

func updateGaugeMetrics(metrics map[string]common.TypeGauge) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	metrics["Alloc"] = common.TypeGauge(rtm.Alloc)
	metrics["BuckHashSys"] = common.TypeGauge(rtm.BuckHashSys)
	metrics["Frees"] = common.TypeGauge(rtm.Frees)
	metrics["GCCPUFraction"] = common.TypeGauge(rtm.GCCPUFraction)
	metrics["GCSys"] = common.TypeGauge(rtm.GCSys)
	metrics["HeapAlloc"] = common.TypeGauge(rtm.HeapAlloc)
	metrics["HeapIdle"] = common.TypeGauge(rtm.HeapIdle)
	metrics["HeapInuse"] = common.TypeGauge(rtm.HeapInuse)
	metrics["HeapObjects"] = common.TypeGauge(rtm.HeapObjects)
	metrics["HeapReleased"] = common.TypeGauge(rtm.HeapReleased)
	metrics["HeapSys"] = common.TypeGauge(rtm.HeapSys)
	metrics["LastGC"] = common.TypeGauge(rtm.LastGC)
	metrics["Lookups"] = common.TypeGauge(rtm.Lookups)
	metrics["MCacheInuse"] = common.TypeGauge(rtm.MCacheInuse)
	metrics["MCacheSys"] = common.TypeGauge(rtm.MCacheSys)
	metrics["MSpanInuse"] = common.TypeGauge(rtm.MSpanInuse)
	metrics["MSpanSys"] = common.TypeGauge(rtm.MSpanSys)
	metrics["Mallocs"] = common.TypeGauge(rtm.Mallocs)
	metrics["NextGC"] = common.TypeGauge(rtm.NextGC)
	metrics["NumForcedGC"] = common.TypeGauge(rtm.NumForcedGC)
	metrics["NumGC"] = common.TypeGauge(rtm.NumGC)
	metrics["OtherSys"] = common.TypeGauge(rtm.OtherSys)
	metrics["PauseTotalNs"] = common.TypeGauge(rtm.PauseTotalNs)
	metrics["StackInuse"] = common.TypeGauge(rtm.StackInuse)
	metrics["StackSys"] = common.TypeGauge(rtm.StackSys)
	metrics["Sys"] = common.TypeGauge(rtm.Sys)
	metrics["TotalAlloc"] = common.TypeGauge(rtm.TotalAlloc)
	metrics["RandomValue"] = common.TypeGauge(rand.Float64())
}
func updateCounterMetrics(metrics map[string]common.TypeCounter) {
	metrics["PollCount"]++
}
func prepareReportGaugeMetrics(metricsGauge map[string]common.TypeGauge) []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metricsGauge {
		temp := float64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "gauge", Value: &temp})
	}
	return dataMetricForReport
}

func prepareReportCounterMetrics(metricsCounter map[string]common.TypeCounter) []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metricsCounter {
		temp := int64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "counter", Delta: &temp})
	}
	return dataMetricForReport
}

func reportMetrics(client *http.Client, config Config, dataMetricForReport []common.Metrics,
	logger common.Loger) {
	for _, metric := range dataMetricForReport {
		err := reportMetricWithRetry(client, config, metric, logger)
		if err != nil {
			logger.Warnf("error report metric. err%s\n ", err)
		}
	}
}

func reportMetricWithRetry(client *http.Client, config Config, metric common.Metrics,
	logger common.Loger) error {
	secondWaitRetry := []int{0, 1, 3, 5}
	var err error
	for _, val := range secondWaitRetry {
		time.Sleep(time.Duration(val) * time.Second)
		err = reportMetric(client, config, metric, logger)
		if !errors.Is(err, syscall.ECONNREFUSED) {
			if err != nil {
				logger.Warnf("error report metric. err%s\n ", err)
			}
			break
		}
	}
	return err
}

func reportMetric(client *http.Client, config Config, metric common.Metrics, logger common.Loger,
) error {
	objMetrics, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("problem with marshal JSON file. err:%w", err)
	}
	url := fmt.Sprintf("http://%s/update/", config.ServerAdderess)

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(objMetrics); err != nil {
		return fmt.Errorf("problem with compress. err:%w", err)
	}
	if err = g.Close(); err != nil {
		return fmt.Errorf("problem with  close compress writer. err:%w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(objMetrics),
	)
	if err != nil {
		logger.Warnf("http.NewRequest. err: %s\n", err)
	}
	const encod = "gzip"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", encod)
	req.Header.Set("Content-Encoding", encod)
	sig, err := common.Sign(objMetrics, config.SignKey)
	if err != nil {
		logger.Warnf("Error sign. err: %s\n", err)
	} else if len(sig) > 0 {
		req.Header.Add("HashSHA256", string(sig))
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Debugf("http.NewRequest.Do err: %s\n", err)
		return fmt.Errorf("http.NewRequest.Do err:%w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Warnf("resp.Body.Close() err: %s\n", err)
		}
	}()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body. err:%w", err)
	}
	return nil
}

func reportCounterMetrics(client *http.Client, config Config, dataMetricForReport []common.Metrics,
	metricsCounter map[string]common.TypeCounter, logger common.Loger) {
	for _, metric := range dataMetricForReport {
		err := reportMetricWithRetry(client, config, metric, logger)
		if err != nil {
			logger.Warnf("error report metric for counter. err%s\n ", err)
			continue
		} else {
			metricsCounter[metric.ID] -= common.TypeCounter(*metric.Delta)
		}
	}
}
func MetricsWatcher(config Config, client *http.Client, logger common.Loger, done chan struct{}) {
	var isBatch = true
	tickerPoolInterval := time.NewTicker(time.Duration(config.PollIntervalSecond) * time.Second)
	tickerReportInterval := time.NewTicker(time.Duration(config.ReportIntervalSecond) * time.Second)
	metricsGauge := make(map[string]common.TypeGauge)
	metricsCounter := make(map[string]common.TypeCounter)
	metricsCounter["PollCount"] = 0
	for {
		select {
		case <-done:
			return
		case <-tickerPoolInterval.C:
			updateGaugeMetrics(metricsGauge)
			updateCounterMetrics(metricsCounter)
		case <-tickerReportInterval.C:
			metricsForReport := prepareReportGaugeMetrics(metricsGauge)
			metricsCounterReport := prepareReportCounterMetrics(metricsCounter)
			if isBatch {
				metricsForReport = append(metricsForReport, metricsCounterReport...)
				err := reportAllMetrics(client, config, metricsForReport, logger)
				if err != nil {
					logger.Warnf("error for send all metrics. err:%s", err)
				} else {
					cleanCounterMetric(metricsCounter)
				}
			} else {
				reportMetrics(client, config, metricsForReport, logger)
				reportCounterMetrics(client, config, metricsCounterReport, metricsCounter, logger)
			}
		}
	}
}

func sendArrayMetric(client *http.Client, config Config, metrics []common.Metrics, logger common.Loger,
) error {
	objMetrics, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("problem with marshal JSON file. err:%w", err)
	}
	url := fmt.Sprintf("http://%s/updates/", config.ServerAdderess)

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(objMetrics); err != nil {
		return fmt.Errorf("problem with compress. err:%w", err)
	}
	if err = g.Close(); err != nil {
		return fmt.Errorf("problem with  close compress writer. err:%w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(objMetrics),
	)
	if err != nil {
		logger.Warnf("http.NewRequest. err: %s\n", err)
	}
	const encod = "gzip"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", encod)
	req.Header.Set("Content-Encoding", encod)
	sig, err := common.Sign(objMetrics, config.SignKey)
	if err != nil {
		logger.Warnf("error sign. err: %s\n", err)
	} else if len(sig) > 0 {
		logger.Debugf("try set HashSHA256 sign: %s", string(sig))
		req.Header.Set("HashSHA256", "aaaa")
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Debugf("http.NewRequest.Do err: %s\n", err)
		return fmt.Errorf("http.NewRequest.Do err:%w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Warnf("resp.Body.Close() err: %s\n", err)
		}
	}()
	return nil
}
func reportAllMetrics(client *http.Client, config Config, dataMetricForReport []common.Metrics,
	logger common.Loger) error {
	err := sendArrayMetric(client, config, dataMetricForReport, logger)
	if err != nil {
		return fmt.Errorf("error reportAllMetrics. err:%w\n ", err)
	}
	return nil
}
func cleanCounterMetric(metricsCounter map[string]common.TypeCounter) {
	for key := range metricsCounter {
		metricsCounter[key] = common.TypeCounter(0)
	}
}
