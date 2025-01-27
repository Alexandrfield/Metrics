package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
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
	"github.com/shirou/gopsutil/v4/mem"
)

func updateGaugeMetrics(metrics *MetricsMap) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	metrics.UpdateGauge("Alloc", common.TypeGauge(rtm.Alloc))
	metrics.UpdateGauge("BuckHashSys", common.TypeGauge(rtm.BuckHashSys))
	metrics.UpdateGauge("Frees", common.TypeGauge(rtm.Frees))
	metrics.UpdateGauge("GCCPUFraction", common.TypeGauge(rtm.GCCPUFraction))
	metrics.UpdateGauge("GCSys", common.TypeGauge(rtm.GCSys))
	metrics.UpdateGauge("HeapAlloc", common.TypeGauge(rtm.HeapAlloc))
	metrics.UpdateGauge("HeapIdle", common.TypeGauge(rtm.HeapIdle))
	metrics.UpdateGauge("HeapInuse", common.TypeGauge(rtm.HeapInuse))
	metrics.UpdateGauge("HeapObjects", common.TypeGauge(rtm.HeapObjects))
	metrics.UpdateGauge("HeapReleased", common.TypeGauge(rtm.HeapReleased))
	metrics.UpdateGauge("HeapSys", common.TypeGauge(rtm.HeapSys))
	metrics.UpdateGauge("LastGC", common.TypeGauge(rtm.LastGC))
	metrics.UpdateGauge("Lookups", common.TypeGauge(rtm.Lookups))
	metrics.UpdateGauge("MCacheInuse", common.TypeGauge(rtm.MCacheInuse))
	metrics.UpdateGauge("MCacheSys", common.TypeGauge(rtm.MCacheSys))
	metrics.UpdateGauge("MSpanInuse", common.TypeGauge(rtm.MSpanInuse))
	metrics.UpdateGauge("MSpanSys", common.TypeGauge(rtm.MSpanSys))
	metrics.UpdateGauge("Mallocs", common.TypeGauge(rtm.Mallocs))
	metrics.UpdateGauge("NextGC", common.TypeGauge(rtm.NextGC))
	metrics.UpdateGauge("NumForcedGC", common.TypeGauge(rtm.NumForcedGC))
	metrics.UpdateGauge("NumGC", common.TypeGauge(rtm.NumGC))
	metrics.UpdateGauge("OtherSys", common.TypeGauge(rtm.OtherSys))
	metrics.UpdateGauge("PauseTotalNs", common.TypeGauge(rtm.PauseTotalNs))
	metrics.UpdateGauge("StackInuse", common.TypeGauge(rtm.StackInuse))
	metrics.UpdateGauge("StackSys", common.TypeGauge(rtm.StackSys))
	metrics.UpdateGauge("Sys", common.TypeGauge(rtm.Sys))
	metrics.UpdateGauge("TotalAlloc", common.TypeGauge(rtm.TotalAlloc))
	metrics.UpdateGauge("RandomValue", common.TypeGauge(rand.Float64()))
}
func SaveAdditionalMetrics(metrics *MetricsMap) {
	v, _ := mem.VirtualMemory()
	metrics.UpdateGauge("TotalMemory", common.TypeGauge(v.TotalMemory))
	metrics.UpdateGauge("FreeMemory", common.TypeGauge(v.FreeMemory))
	metrics.UpdateGauge("CPUutilization1", common.TypeGauge(rand.Float64()))
}
func updateCounterMetrics(metrics *MetricsMap) {
	metrics.UpdateCounter("PollCount", common.TypeCounter(1))
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
		req.Header.Add("HashSHA256", hex.EncodeToString(sig))
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
	metrics *MetricsMap, logger common.Loger) {
	for _, metric := range dataMetricForReport {
		err := reportMetricWithRetry(client, config, metric, logger)
		if err != nil {
			logger.Warnf("error report metric for counter. err%s\n ", err)
			continue
		} else {
			metrics.UpdateCounter(metric.ID, (-1)*common.TypeCounter(*metric.Delta))
		}
	}
}
func MetricsWatcher(config Config, client *http.Client, logger common.Loger, done chan struct{}) {
	var isBatch = true
	tickerPoolInterval := time.NewTicker(time.Duration(config.PollIntervalSecond) * time.Second)
	tickerReportInterval := time.NewTicker(time.Duration(config.ReportIntervalSecond) * time.Second)
	metrics := MetricsMap{}
	metrics.Initializate()
	for {
		select {
		case <-done:
			return
		case <-tickerPoolInterval.C:
			updateGaugeMetrics(metrics)
			updateCounterMetrics(metrics)
		case <-tickerReportInterval.C:
			metricsForReport := metrics.PrepareReportGaugeMetrics()
			metricsCounterReport := metrics.PrepareReportCounterMetrics()
			if isBatch {
				metricsForReport = append(metricsForReport, metricsCounterReport...)
				err := reportAllMetrics(client, config, metricsForReport, logger)
				if err != nil {
					logger.Warnf("error for send all metrics. err:%s", err)
				} else {
					//cleanCounterMetric(metricsCounter)
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
		logger.Debugf("try set HashSHA256 sign: %s", hex.EncodeToString(sig))
		req.Header.Set("HashSHA256", hex.EncodeToString(sig))
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
