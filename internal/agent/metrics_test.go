package agent

import (
	"testing"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAdditionalMetrics(t *testing.T) {
	var listMetricsName = []string{"TotalMemory", "FreeMemory", "CPUutilization1"}
	metrics := MetricsMap{}
	metrics.Initializate()
	for _, val := range listMetricsName {
		metrics.UpdateGauge(val, -1)
	}
	updateAdditionalMetrics(&metrics)
	for _, val := range listMetricsName {
		v := metrics.GetGauge(val)
		if v == -1 {
			t.Errorf("Error key:%s; value:%v\n", val, v)
		}
	}
}

func TestAdditionalMetricsWatcher(t *testing.T) {
	var listMetricsName = []string{"TotalMemory", "FreeMemory", "CPUutilization1"}
	metrics := MetricsMap{}
	metrics.Initializate()
	for _, val := range listMetricsName {
		metrics.UpdateGauge(val, -1)
	}
	done := make(chan struct{})
	config := Config{PollIntervalSecond: 1}

	go AdditionalMetricsWatcher(config, &metrics, done)
	time.Sleep(2 * time.Second)
	close(done)

	for _, val := range listMetricsName {
		v := metrics.GetGauge(val)
		if v == -1 {
			t.Errorf("Error key:%s; value:%v\n", val, v)
		}
	}
}
func TestUpdateGaugeMetrics(t *testing.T) {
	var listMetricsName = []string{"GCCPUFraction",
		"GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
		"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue"}
	metrics := MetricsMap{}
	metrics.Initializate()
	for _, val := range listMetricsName {
		metrics.UpdateGauge(val, -1)
	}
	updateGaugeMetrics(&metrics)
	for _, val := range listMetricsName {
		v := metrics.GetGauge(val)
		if v == -1 {
			t.Errorf("Error key:%s; value:%v\n", val, v)
		}
	}
}

func TestUpdateCounterMetrics(t *testing.T) {
	var listMetricsName = []string{"PollCount"}
	metrics := MetricsMap{}
	metrics.Initializate()
	for _, val := range listMetricsName {
		metrics.UpdateCounter(val, -1)
	}
	updateCounterMetrics(&metrics)
	for _, val := range listMetricsName {
		v := metrics.GetCounter(val)
		if v == -1 {
			t.Errorf("Error key:%s; value:%v\n", val, v)
		}
	}
}

func TestPrepareReportGaugeMetrics(t *testing.T) {
	metrics := MetricsMap{}
	metrics.Initializate()
	listIds := []string{"Alloc", "GCCPUFraction"}

	expected := make([]common.Metrics, 0)
	for _, value := range listIds {
		temp := 9.1
		expected = append(expected, common.Metrics{ID: value, MType: "gauge", Value: &temp})
		metrics.UpdateGauge(value, common.TypeGauge(temp))
	}
	actual := metrics.PrepareReportGaugeMetrics()
	assert.ElementsMatch(t, actual, expected)
}
func TestPrepareReportCounterMetrics(t *testing.T) {
	metrics := MetricsMap{}
	metrics.Initializate()
	listIds := []string{"AllocCounter", "GCCPUFractionCounter"}

	expected := make([]common.Metrics, 0)
	for _, value := range listIds {
		temp := int64(10)
		expected = append(expected, common.Metrics{ID: value, MType: "counter", Delta: &temp})
		metrics.UpdateCounter(value, common.TypeCounter(temp))
	}
	actual := metrics.PrepareReportCounterMetrics()
	assert.ElementsMatch(t, actual, expected)
}
func TestFixSusceeseSavedCounterMetric(t *testing.T) {
	metrics := MetricsMap{}
	metrics.Initializate()
	var temp1 int64 = 4
	var temp2 int64 = 56
	testData := []common.Metrics{
		{ID: "test1", MType: "counter", Delta: &temp1},
		{ID: "test2", MType: "counter", Delta: &temp2},
	}
	for _, val := range testData {
		metrics.UpdateCounter(val.ID, common.TypeCounter(*val.Delta*2))
	}
	fixSusceeseSavedCounterMetric(&metrics, testData)

	for _, val := range testData {
		actual := int64(metrics.GetCounter(val.ID))
		if *val.Delta != actual {
			t.Errorf("Error value for ID:%s. expected:%d, actual:%d", val.ID, *val.Delta, actual)
		}
	}
}
