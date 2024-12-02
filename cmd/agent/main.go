package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Alexandrfield/Metrics/internal/storage"
)

var globalPollIntervalSecond int = 2
var globalReportIntervalSecond int = 10
var globalServerAdderess string = "127.0.0.1:8080"

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
func reportMetrics(client *http.Client, metricsGauge map[string]storage.TypeGauge, metricsCounter map[string]storage.TypeCounter) {

	dataMetricForReport := make([]string, 0)
	for key, value := range metricsGauge {
		dataMetricForReport = append(dataMetricForReport, fmt.Sprintf("gauge/%s/%s", key, value))
	}
	for key, value := range metricsCounter {
		dataMetricForReport = append(dataMetricForReport, fmt.Sprintf("counter/%s/%s", key, value))
	}
	go func() {
		for _, metric := range dataMetricForReport {
			reportMetric(client, metric)
		}
	}()
}
func reportMetric(client *http.Client, metric string) {
	url := fmt.Sprintf("http://%s/update/%s", globalServerAdderess, metric)
	req, err := http.NewRequest(
		http.MethodPost, url, nil,
	)
	if err != nil {
		fmt.Printf("http.NewRequest. err: %s\n", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http.NewRequest.Do err: %s\n", err)
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Printf("Error reading body. err%s\n ", err)
	}
}
func metricsWatcher(client *http.Client, done chan struct{}) {
	tickerPoolInterval := time.NewTicker(time.Duration(globalPollIntervalSecond) * time.Second)
	tickerReportInterval := time.NewTicker(time.Duration(globalReportIntervalSecond) * time.Second)
	metricsGauge := make(map[string]storage.TypeGauge)
	metricsCounter := make(map[string]storage.TypeCounter)
	metricsCounter["PollCount"] = 0
	for {
		select {
		case <-done:
			return
		case <-tickerPoolInterval.C:
			updateGaugeMetrics(metricsGauge)
			metricsCounter["PollCount"] = metricsCounter["PollCount"] + 1
		case <-tickerReportInterval.C:
			reportMetrics(client, metricsGauge, metricsCounter)
		}
	}
}
func main() {

	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}

	done := make(chan struct{})
	go metricsWatcher(&client, done)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-osSignals
	close(done)
	time.Sleep(1 * time.Second)
}
