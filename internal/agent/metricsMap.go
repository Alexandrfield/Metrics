package agent

import (
	"sync"

	"github.com/Alexandrfield/Metrics/internal/common"
)

// MetricsMap Object for save all type metric.
type MetricsMap struct {
	metricsGauge        map[string]common.TypeGauge
	metricsCounter      map[string]common.TypeCounter
	mutexMetricsGauge   sync.RWMutex
	mutexMetricsCounter sync.RWMutex
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

// PrepareReportGaugeMetrics The function collects that GaugeMetrics
// that stored metrics and prepares them for sending as an array.
func (metr *MetricsMap) PrepareReportGaugeMetrics() []common.Metrics {
	metr.mutexMetricsGauge.RLock()
	defer metr.mutexMetricsGauge.RUnlock()
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metr.metricsGauge {
		temp := float64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "gauge", Value: &temp})
	}
	return dataMetricForReport
}

// PrepareReportCounterMetrics The function collects all currently
// stored metrics and prepares them for sending as an array.
func (metr *MetricsMap) PrepareReportCounterMetrics() []common.Metrics {
	metr.mutexMetricsCounter.RLock()
	defer metr.mutexMetricsCounter.RUnlock()
	dataMetricForReport := make([]common.Metrics, 0)
	for key, value := range metr.metricsCounter {
		temp := int64(value)
		dataMetricForReport = append(dataMetricForReport, common.Metrics{ID: key, MType: "counter", Delta: &temp})
	}
	return dataMetricForReport
}
