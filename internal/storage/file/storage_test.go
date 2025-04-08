package filestorage

import (
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGauge(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	dataTest := []struct {
		name  string
		value common.TypeGauge
	}{
		{
			name:  "testGauge1",
			value: common.TypeGauge(6.1),
		},
		{
			name:  "testGauge1",
			value: common.TypeGauge(6.1),
		},
		{
			name:  "testGauge2",
			value: common.TypeGauge(7.89),
		},
	}
	for _, val := range dataTest {
		stor.AddGauge(val.name, val.value)
	}

	for _, val := range dataTest {
		v, ok := stor.GaugeData[val.name]
		if !ok || v != val.value {
			t.Errorf("not save gauge metric. name:%s; actual:%v; expected:%v; ok:%t", val.name, v, val.value, ok)
		}
	}
}

func TestAddCounter(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	dataTest := []struct {
		name  string
		value common.TypeCounter
	}{
		{
			name:  "testGauge1",
			value: common.TypeCounter(4),
		},
		{
			name:  "testGauge2",
			value: common.TypeCounter(49),
		},
	}
	for _, val := range dataTest {
		stor.AddCounter(val.name, val.value)
		stor.AddCounter(val.name, val.value)
	}

	for _, val := range dataTest {
		v, ok := stor.CounterData[val.name]
		if !ok || v != val.value*2 {
			t.Errorf("not save counter metric. name:%s; actual:%v; expected:%v; ok:%t", val.name, v, val.value*2, ok)
		}
	}
}

func TestGetGauge(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	dataTest := []struct {
		name  string
		value common.TypeGauge
	}{
		{
			name:  "testGauge1",
			value: common.TypeGauge(6.1),
		},
		{
			name:  "testGauge2",
			value: common.TypeGauge(7.89),
		},
	}
	for _, val := range dataTest {
		stor.GaugeData[val.name] = val.value
	}

	for _, val := range dataTest {
		actual, err := stor.GetGauge(val.name)
		require.NoError(t, err)
		assert.Equal(t, actual, val.value)
	}
}
func TestGetCounter(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	dataTest := []struct {
		name  string
		value common.TypeCounter
	}{
		{
			name:  "testGauge1",
			value: common.TypeCounter(4),
		},
		{
			name:  "testGauge2",
			value: common.TypeCounter(49),
		},
	}
	for _, val := range dataTest {
		stor.CounterData[val.name] = val.value
	}

	for _, val := range dataTest {
		actual, err := stor.GetCounter(val.name)
		require.NoError(t, err)
		assert.Equal(t, actual, val.value)
	}
}

func TestPingDatabase(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	actual := stor.PingDatabase()
	if actual {
		t.Errorf("Wrong result Ping databese. actual:%t; ecpected:false", actual)
	}
}

func TestGetAllMetricName(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	nameGaugeMetrics := []string{"testGauge1", "testGauge2", "testGauge3"}
	nameCounterMetrics := []string{"testCounter1", "testCounter2", "testCounter3"}
	for _, val := range nameGaugeMetrics {
		stor.GaugeData[val] = 3.0
	}
	for _, val := range nameCounterMetrics {
		stor.CounterData[val] = 4
	}

	actualNameGauge, actualNameCounter := stor.GetAllMetricName()
	assert.ElementsMatch(t, actualNameGauge, nameGaugeMetrics)
	assert.ElementsMatch(t, actualNameCounter, nameCounterMetrics)
}

func TestAddMetrics(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []float64{6.7, 7.89}
	metricsCounterName := []string{"testCounter1", "testCounter2"}
	metricsCounterval := []int64{56, 12}
	testMetrics := []common.Metrics{
		{ID: metricsGaugeName[0], MType: "gauge", Value: &metricsGaugeVal[0]},
		{ID: metricsGaugeName[1], MType: "gauge", Value: &metricsGaugeVal[1]},
		{ID: metricsCounterName[0], MType: "counter", Delta: &metricsCounterval[0]},
		{ID: metricsCounterName[1], MType: "counter", Delta: &metricsCounterval[1]},
		{ID: metricsCounterName[0], MType: "counter", Delta: &metricsCounterval[0]},
		{ID: metricsCounterName[1], MType: "counter", Delta: &metricsCounterval[1]},
	}

	err := stor.AddMetrics(testMetrics)
	require.NoError(t, err)
	for i, val := range metricsGaugeName {
		v, ok := stor.GaugeData[val]
		expect := common.TypeGauge(metricsGaugeVal[i])
		if !ok || v != expect {
			t.Errorf("not save gauge metric. name:%s; actual:%v; expected:%v; ok:%t", val, v, expect, ok)
		}
	}
	for i, val := range metricsCounterName {
		v, ok := stor.CounterData[val]
		expect := common.TypeCounter(metricsCounterval[i] * 2)
		if !ok || v != expect {
			t.Errorf("not save gauge metric. name:%s; actual:%v; expected:%v; ok:%t", val, v, expect, ok)
		}
	}
}
func TestAddMetricsNegativ(t *testing.T) {
	stor := NewMemFileStorage(&common.FakeLogger{})
	metricsGaugeVal := 45.78
	testMetrics := []common.Metrics{
		{ID: "test", MType: "blblabla", Value: &metricsGaugeVal},
	}

	err := stor.AddMetrics(testMetrics)
	if err == nil {
		t.Errorf("no detect error metric. actual err nil; ecpected: !nil")
	}
}
