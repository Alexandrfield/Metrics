package agent

import (
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGaugeMetrics(t *testing.T) {
	var listMetricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction",
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
