package agent

import (
	"sync"

	"github.com/Alexandrfield/Metrics/internal/common"
)

type MetricsMap struct {
	mutexMetricsGauge   sync.RWMutex
	metricsGauge        map[string]common.TypeGauge
	mutexMetricsCounter sync.RWMutex
	metricsCounter      map[string]common.TypeCounter
}

func (metr *MetricsMap) Initializate() {
	metr.metricsGauge = make(map[string]common.TypeGauge)
	metr.metricsCounter = make(map[string]common.TypeCounter)
}
func (metr *MetricsMap) UpdateGauge(name string, value common.TypeGauge) {
	metr.mutexMetricsGauge.Lock()
	defer metr.mutexMetricsGauge.Unlock()
	metr.metricsGauge[name] = value
}
func (metr *MetricsMap) UpdateCounter(name string, value common.TypeCounter) {
	metr.mutexMetricsCounter.Lock()
	defer metr.mutexMetricsCounter.Unlock()
	val, ok := metr.metricsCounter[name]
	if !ok {
		val = 0
	}
	metr.metricsCounter[name] = value + val
}
func (metr *MetricsMap) GetGauge(name string) common.TypeGauge {
	metr.mutexMetricsGauge.RLock()
	defer metr.mutexMetricsGauge.RUnlock()
	return metr.metricsGauge[name]
}
func (metr *MetricsMap) GetCounter(name string) common.TypeCounter {
	metr.mutexMetricsCounter.RLock()
	defer metr.mutexMetricsCounter.RUnlock()
	return metr.metricsCounter[name]
}
func (metr *MetricsMap) PrepareReportGaugeMetrics() []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metr.metricsGauge {
		temp := float64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "gauge", Value: &temp})
	}
	return dataMetricForReport
}

func (metr *MetricsMap) PrepareReportCounterMetrics() []common.Metrics {
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metr.metricsCounter {
		temp := int64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "counter", Delta: &temp})
	}
	return dataMetricForReport
}
