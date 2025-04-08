package agent

import (
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
)

func TestSetMetricsMapUpdateGauge(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeGauge(64.0)
	metrMap.UpdateGauge(testMetric, testValue)
	val, ok := metrMap.metricsGauge[testMetric]
	if !ok {
		t.Errorf("test metric name:%s, val:%f; is not save in MetricsMap!", testMetric, testValue)
		return
	}
	if val != testValue {
		t.Errorf("test metric name:%s, expect:%f, actual:%f", testMetric, testValue, val)
		return
	}
}
func TestUpdateMetricsMapUpdateGauge(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeGauge(64.0)
	metrMap.UpdateGauge(testMetric, testValue)
	metrMap.UpdateGauge(testMetric, testValue)
	val, ok := metrMap.metricsGauge[testMetric]
	if !ok {
		t.Errorf("test metric name:%s, val:%f; is not save in MetricsMap!", testMetric, testValue)
		return
	}
	updatedValue := testValue
	if val != updatedValue {
		t.Errorf("test metric name:%s, expect:%f, actual:%f", testMetric, testValue, val)
		return
	}
}
func TestSetMetricsMapUpdateCounter(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeCounter(31)
	metrMap.UpdateCounter(testMetric, testValue)
	val, ok := metrMap.metricsCounter[testMetric]
	if !ok {
		t.Errorf("test metric name:%s, val:%d; is not save in MetricsMap!", testMetric, testValue)
		return
	}
	if val != testValue {
		t.Errorf("test metric name:%s, expect:%d, actual:%d", testMetric, testValue, val)
		return
	}
}
func TestUpdateMetricsMapUpdateCounter(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeCounter(31)
	metrMap.UpdateCounter(testMetric, testValue)
	metrMap.UpdateCounter(testMetric, testValue)
	val, ok := metrMap.metricsCounter[testMetric]
	if !ok {
		t.Errorf("test metric name:%s, val:%d; is not save in MetricsMap!", testMetric, testValue)
		return
	}
	updatedValue := testValue + testValue
	if val != updatedValue {
		t.Errorf("test metric name:%s, expect:%d, actual:%d", testMetric, updatedValue, val)
		return
	}
}
func TestMetricsMapGetGauge(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeGauge(64.0)
	metrMap.metricsGauge[testMetric] = testValue
	val := metrMap.GetGauge(testMetric)
	if val != testValue {
		t.Errorf("test metric name:%s, expect:%f, actual:%f", testMetric, testValue, val)
		return
	}
}
func TestMetricsMapGetCounter(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()
	testMetric := "TestMetric"
	testValue := common.TypeCounter(29)
	metrMap.metricsCounter[testMetric] = testValue

	val := metrMap.GetCounter(testMetric)
	if val != testValue {
		t.Errorf("test metric name:%s, expect:%d, actual:%d", testMetric, testValue, val)
		return
	}
}
func TestMetricsMapPrepareReportGaugeMetrics(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()

	val1 := float64(64.0)
	val2 := float64(31.4)
	val3 := float64(11.22)
	testData := []common.Metrics{
		{ID: "test1", MType: "gauge", Value: &val1},
		{ID: "test2", MType: "gauge", Value: &val2},
		{ID: "test3", MType: "gauge", Value: &val3},
	}

	for _, val := range testData {
		metrMap.UpdateGauge(val.ID, common.TypeGauge(*val.Value))
	}
	actual := metrMap.PrepareReportGaugeMetrics()
	for _, expectVal := range testData {
		res := false
		for _, actualVal := range actual {
			if actualVal.ID == expectVal.ID && *actualVal.Value == *expectVal.Value {
				res = true
				break
			}
		}
		if !res {
			t.Errorf("Cant find metric: id:%s,value%f in actual. Len(actual):%d", expectVal.ID, *expectVal.Value, len(actual))
			return
		}
	}
}

func TestMetricsMapPrepareReportCounterMetrics(t *testing.T) {
	metrMap := MetricsMap{}
	metrMap.Initializate()

	val1 := int64(56)
	val2 := int64(100)
	val3 := int64(5)
	testData := []common.Metrics{
		{ID: "test1", MType: "counter", Delta: &val1},
		{ID: "test2", MType: "counter", Delta: &val2},
		{ID: "test3", MType: "counter", Delta: &val3},
	}

	for _, val := range testData {
		metrMap.UpdateCounter(val.ID, common.TypeCounter(*val.Delta))
	}
	actual := metrMap.PrepareReportCounterMetrics()
	for _, expectVal := range testData {
		res := false
		for _, actualVal := range actual {
			if actualVal.ID == expectVal.ID && *actualVal.Delta == *expectVal.Delta {
				res = true
				break
			}
		}
		if !res {
			t.Errorf("Cant find metric: id:%s,value%d in actual. Len(actual):%d", expectVal.ID, *expectVal.Delta, len(actual))
			return
		}
	}
}
