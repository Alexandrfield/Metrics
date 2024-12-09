package agent

import (
	"fmt"
	"testing"

	"github.com/Alexandrfield/Metrics/internal/storage"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGaugeMetrics(t *testing.T) {
	var listMetricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction",
		"GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
		"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue"}
	metricsGauge := make(map[string]storage.TypeGauge)
	for _, val := range listMetricsName {
		metricsGauge[val] = -1
	}
	updateGaugeMetrics(metricsGauge)
	for _, val := range listMetricsName {
		v, ok := metricsGauge[val]
		if !ok || v == -1 {
			t.Errorf("Error key:%s;ok?:%v; value:%v\n", val, ok, v)
		}
		if ok {
			delete(metricsGauge, val)
		}
	}
	for key, val := range metricsGauge {
		t.Errorf("Error unexpected key :%s; value:%v;\n", key, val)
	}
}

func TestUpdateCounterMetrics(t *testing.T) {
	var listMetricsName = []string{"PollCount"}
	metricsCounter := make(map[string]storage.TypeCounter)
	for _, val := range listMetricsName {
		metricsCounter[val] = -1
	}
	updateCounterMetrics(metricsCounter)
	for _, val := range listMetricsName {
		v, ok := metricsCounter[val]
		if !ok || v == -1 {
			t.Errorf("Error key:%s;ok?:%v; value:%v\n", val, ok, v)
		}
		if ok {
			delete(metricsCounter, val)
		}
	}
	for key, val := range metricsCounter {
		t.Errorf("Error unexpected key :%s; value:%v;\n", key, val)
	}
}

func TestPrepareReportMetrics(t *testing.T) {
	metricsCounter := make(map[string]storage.TypeCounter)
	metricsCounter["PollCount"] = 100
	metricsCounter["test"] = -1
	metricsGauge := make(map[string]storage.TypeGauge)
	metricsGauge["Alloc"] = 9.1
	metricsGauge["GCCPUFraction"] = 10.43
	serverAdderess := "127.0.0.1:8080"

	var expected = []string{
		fmt.Sprintf("http://%s/update/gauge/Alloc/9.1", serverAdderess),
		fmt.Sprintf("http://%s/update/gauge/GCCPUFraction/10.43", serverAdderess),
		fmt.Sprintf("http://%s/update/counter/PollCount/100", serverAdderess),
		fmt.Sprintf("http://%s/update/counter/test/-1", serverAdderess),
	}
	actual := prepareReportMetrics(serverAdderess, metricsGauge, metricsCounter)
	assert.ElementsMatch(t, actual, expected)
}
