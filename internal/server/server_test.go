package server

import (
	"fmt"
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	mockServ "github.com/Alexandrfield/Metrics/internal/server/mock"
	"github.com/Alexandrfield/Metrics/internal/storage"
	mock "github.com/Alexandrfield/Metrics/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponceWriter := mockServ.NewMockResponseWriter(ctrl)
	b := []byte{0x01, 0x02}
	mockResponceWriter.EXPECT().Write(b).Return(1, nil).Times(1)
	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := loggingResponseWriter{ResponseWriter: mockResponceWriter, responseData: responseData}
	n, err := lw.Write(b)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
}
func TestWriteErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponceWriter := mockServ.NewMockResponseWriter(ctrl)
	b := []byte{0x01, 0x02}
	errT := fmt.Errorf("loggingResponseWriter errTest")
	mockResponceWriter.EXPECT().Write(b).Return(1, errT).Times(1)
	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := loggingResponseWriter{ResponseWriter: mockResponceWriter, responseData: responseData}
	n, err := lw.Write(b)
	if err == nil {
		t.Errorf("expected:nil; actual:%s", err)
	}
	assert.Equal(t, 1, n)
}
func TestWriteHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponceWriter := mockServ.NewMockResponseWriter(ctrl)
	statusCodeTest := 3
	mockResponceWriter.EXPECT().WriteHeader(statusCodeTest).Times(1)
	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := loggingResponseWriter{ResponseWriter: mockResponceWriter, responseData: responseData}
	lw.WriteHeader(statusCodeTest)
	if lw.responseData.status != statusCodeTest {
		t.Errorf("expected:%d; actual:%d", statusCodeTest, lw.responseData.status)
	}
}

func TestSetCounterValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeCounter(65)

	mockBasicStorage.EXPECT().AddCounter(metricName, metricValue).Return(nil).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	err := testRep.SetCounterValue(metricName, metricValue)
	require.NoError(t, err)
}
func TestSetCounterValueNilStor(t *testing.T) {
	metricName := "testName"
	metricValue := common.TypeCounter(65)

	testRep := CreateMetricRepository(nil, &common.FakeLogger{})

	err := testRep.SetCounterValue(metricName, metricValue)
	if err == nil {
		t.Errorf("Expected error!")
	}
}
func TestSetCounterValueErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeCounter(65)
	errT := fmt.Errorf("problem SetCounterValue. err test")
	mockBasicStorage.EXPECT().AddCounter(metricName, metricValue).Return(errT).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	err := testRep.SetCounterValue(metricName, metricValue)
	if err == nil {
		t.Errorf("Expected error! actual er:nil")
	}
}

func TestSetGaugeValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeGauge(6.78)

	mockBasicStorage.EXPECT().AddGauge(metricName, metricValue).Return(nil).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	err := testRep.SetGaugeValue(metricName, metricValue)
	require.NoError(t, err)
}

func TestSetGaugeValueErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeGauge(6.78)
	errT := fmt.Errorf("problem SetGaugeValue. err test")
	mockBasicStorage.EXPECT().AddGauge(metricName, metricValue).Return(errT).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	err := testRep.SetGaugeValue(metricName, metricValue)
	if err == nil {
		t.Errorf("Expected error! actual er:nil")
	}
}
func TestSetGaugeValueNilStor(t *testing.T) {
	metricName := "testName"
	metricValue := common.TypeGauge(65)

	testRep := CreateMetricRepository(nil, &common.FakeLogger{})

	err := testRep.SetGaugeValue(metricName, metricValue)
	if err == nil {
		t.Errorf("Expected error!")
	}
}
func TestGetCounterValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeCounter(45)

	mockBasicStorage.EXPECT().GetCounter(metricName).Return(metricValue, nil).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	actual, err := testRep.GetCounterValue(metricName)
	require.NoError(t, err)
	assert.Equal(t, metricValue, actual)
}

func TestGetCounterValueErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeCounter(45)

	errT := fmt.Errorf("problem GetCounterValue. err test")
	mockBasicStorage.EXPECT().GetCounter(metricName).Return(metricValue, errT).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	_, err := testRep.GetCounterValue(metricName)
	if err == nil {
		t.Error("Exprcted error. actual err: nil")
	}
}

func TestGetCounterValueNilStor(t *testing.T) {
	metricName := "testName"

	testRep := CreateMetricRepository(nil, &common.FakeLogger{})

	actual, err := testRep.GetCounterValue(metricName)
	if err == nil {
		t.Errorf("Expected error!")
	}
	assert.Equal(t, common.TypeCounter(0), actual)
}

func TestGetGaugeValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeGauge(4.5)

	mockBasicStorage.EXPECT().GetGauge(metricName).Return(metricValue, nil).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	actual, err := testRep.GetGaugeValue(metricName)
	require.NoError(t, err)
	assert.Equal(t, metricValue, actual)
}
func TestGetGaugeValueErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricName := "testName"
	metricValue := common.TypeGauge(4.5)

	errT := fmt.Errorf("problem GetGaugeValue test error")
	mockBasicStorage.EXPECT().GetGauge(metricName).Return(metricValue, errT).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	_, err := testRep.GetGaugeValue(metricName)
	if err == nil {
		t.Error("Exprcted error. actual err: nil")
	}
}

func TestGetGaugeValueNilStore(t *testing.T) {
	metricName := "testName"
	testRep := CreateMetricRepository(nil, &common.FakeLogger{})

	actual, err := testRep.GetGaugeValue(metricName)
	if err == nil {
		t.Errorf("Expected error!")
	}
	assert.Equal(t, common.TypeGauge(0), actual)
}
func TestGetAllValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []common.TypeGauge{6.7, 7.89}
	metricsCounterName := []string{"testCounter1", "testCounter2"}
	metricsCounterval := []common.TypeCounter{56, 12}
	mockBasicStorage.EXPECT().GetAllMetricName().Return(metricsGaugeName, metricsCounterName).Times(1)
	mockBasicStorage.EXPECT().GetGauge(metricsGaugeName[0]).Return(metricsGaugeVal[0], nil).Times(1)
	mockBasicStorage.EXPECT().GetGauge(metricsGaugeName[1]).Return(metricsGaugeVal[1], nil).Times(1)
	mockBasicStorage.EXPECT().GetCounter(metricsCounterName[0]).Return(metricsCounterval[0], nil).Times(1)
	mockBasicStorage.EXPECT().GetCounter(metricsCounterName[1]).Return(metricsCounterval[1], nil).Times(1)

	expected := []string{
		fmt.Sprintf("name:%s; value:%v;\n", metricsGaugeName[0], metricsGaugeVal[0]),
		fmt.Sprintf("name:%s; value:%v;\n", metricsGaugeName[1], metricsGaugeVal[1]),
		fmt.Sprintf("name:%s; value:%v;\n", metricsCounterName[0], metricsCounterval[0]),
		fmt.Sprintf("name:%s; value:%v;\n", metricsCounterName[1], metricsCounterval[1]),
	}

	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	actual, err := testRep.GetAllValue()
	require.NoError(t, err)
	assert.ElementsMatch(t, actual, expected)
}

func TestGetAllValueNilStore(t *testing.T) {
	testRep := CreateMetricRepository(nil, &common.FakeLogger{})
	actual, err := testRep.GetAllValue()
	if err == nil {
		t.Errorf("Expected error!")
	}
	if len(actual) != 0 {
		t.Errorf("Expected len() = 0!")
	}
}
func TestPingDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)
	mockBasicStorage.EXPECT().PingDatabase().Return(true).Times(1)
	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})
	actual := testRep.PingDatabase()
	if !actual {
		t.Errorf("problem ping storage. expected:true; actual:%t", actual)
	}
}

func TestAddMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []float64{6.7, 7.89}
	metricsCounterName := []string{"testCounter1", "testCounter2"}
	metricsCounterval := []int64{56, 12}

	testMetrics := []common.Metrics{
		{ID: metricsGaugeName[0], MType: "gauge", Value: &metricsGaugeVal[0]},
		{ID: metricsGaugeName[1], MType: "gauge", Value: &metricsGaugeVal[1]},
		{ID: metricsCounterName[0], MType: "counter", Delta: &metricsCounterval[0]},
		{ID: metricsCounterName[1], MType: "counter", Delta: &metricsCounterval[1]},
	}

	mockBasicStorage.EXPECT().AddMetrics(testMetrics).Return(nil).Times(1)

	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})

	err := testRep.AddMetrics(testMetrics)
	require.NoError(t, err)
}
func TestAddMetricsNegativ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockBasicStorage(ctrl)

	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []float64{6.7, 7.89}

	testMetrics := []common.Metrics{
		{ID: metricsGaugeName[0], MType: "gauge", Value: &metricsGaugeVal[0]},
		{ID: metricsGaugeName[1], MType: "gauge", Value: &metricsGaugeVal[1]},
	}

	errT := fmt.Errorf("problem wit add metrics. err test")
	mockBasicStorage.EXPECT().AddMetrics(testMetrics).Return(errT).Times(1)

	testRep := CreateMetricRepository(mockBasicStorage, &common.FakeLogger{})
	err := testRep.AddMetrics(testMetrics)
	if err == nil {
		t.Errorf("expected error. actual err: nil")
	}
}

func TestHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponceWriter := mockServ.NewMockResponseWriter(ctrl)
	mockResponceWriter.EXPECT().Header().Return(nil)
	r := mockResponceWriter.Header()
	if r != nil {
		t.Errorf("error call header")
	}
}
func TestWithLoging(t *testing.T) {
	logger := &common.FakeLogger{}
	storageConfig := storage.Config{}
	done := make(chan struct{})
	stor := storage.CreateMemStorage(storageConfig, logger, done)
	metricRep := CreateMetricRepository(stor, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	conf := Config{}
	_ = WithLogging(&common.FakeLogger{}, &conf, servHandler.Ping)
}
