package requesthandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	mock "github.com/Alexandrfield/Metrics/internal/requestHandler/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

var signKey = []byte{0x01, 0x02}

func TestParseURL(t *testing.T) {
	testG1V := 32.5
	testG2V := float64(24)
	testC1V := int64(55)
	testC2V := int64(12)
	testData := []struct {
		name   string
		url    string
		metr   common.Metrics
		status int
	}{
		{
			name:   "positiv gauge 1",
			url:    "localhost/update/gauge/testG1/32.5",
			metr:   common.Metrics{ID: "testG1", Value: &testG1V, MType: "gauge"},
			status: http.StatusOK,
		},
		{
			name:   "positiv gauge 2",
			url:    "localhost/update/gauge/testG2/24",
			metr:   common.Metrics{ID: "testG2", Value: &testG2V, MType: "gauge"},
			status: http.StatusOK,
		},
		{
			name:   "positiv counter 1",
			url:    "localhost/update/counter/testC1/55",
			metr:   common.Metrics{ID: "testC1", Delta: &testC1V, MType: "counter"},
			status: http.StatusOK,
		},
		{
			name:   "positiv counter 2",
			url:    "localhost/update/counter/testC2/12",
			metr:   common.Metrics{ID: "testC2", Delta: &testC2V, MType: "counter"},
			status: http.StatusOK,
		},
		{
			name:   "negativ counter 2",
			url:    "http://localhost/update/counter/testC2/12",
			metr:   common.Metrics{ID: "testC2", Delta: &testC2V, MType: "counter"},
			status: http.StatusNotFound,
		},
		{
			name:   "positiv counter 3",
			url:    "localhost/value/counter/testC3",
			metr:   common.Metrics{ID: "testC3", MType: "counter"},
			status: http.StatusOK,
		},
		{
			name:   "positiv counter 4",
			url:    "http://localhost/value/counter/testC4",
			metr:   common.Metrics{ID: "testC4", MType: "counter"},
			status: http.StatusNotFound,
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {

			metr, st := parseURL(test.url, &common.FakeLogger{})
			if st != test.status {
				t.Errorf("problem parse url. status actual%d; expected:%d", st, test.status)
				return
			}
			if st != http.StatusOK {
				return
			}
			if metr.ID != test.metr.ID {
				t.Errorf("problem parse url. ID actual%s; expected:%s", metr.ID, test.metr.ID)
				return
			}
			if metr.MType != test.metr.MType {
				t.Errorf("problem parse url. MType actual%s; expected:%s", metr.MType, test.metr.MType)
				return
			}
			if test.metr.Value != nil {
				if *metr.Value != *test.metr.Value {
					t.Errorf("problem parse url. Value actual%f; expected:%f", *metr.Value, *test.metr.Value)
					return
				}
			}
			if test.metr.Delta != nil {
				if *metr.Delta != *test.metr.Delta {
					t.Errorf("problem parse url. Delta actual%d; expected:%d", *metr.Delta, *test.metr.Delta)
					return
				}
			}
		})
	}
}

func TestMetricServerPing(t *testing.T) {
	testData := []struct {
		name        string
		checkResult bool
		status      int
	}{
		{
			name:        "positiv check",
			checkResult: true,
			status:      http.StatusOK,
		},
		{
			name:        "negativ check",
			checkResult: false,
			status:      http.StatusInternalServerError,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBasicStorage := mock.NewMockMetricsStorage(ctrl)
			mockBasicStorage.EXPECT().PingDatabase().Return(true).Times(1)
			w := httptest.NewRecorder()
			mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}
			expectdStatus := http.StatusOK

			mServ.Ping(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, expectdStatus, res.StatusCode)
		})
	}
}

func TestUpdateValue(t *testing.T) {
	testG1V := 54.7
	testC2V := int64(5)
	testData := []struct {
		name string
		metr common.Metrics
		err  error
	}{
		{
			name: "positiv check",
			metr: common.Metrics{ID: "testG1", Value: &testG1V, MType: "gauge"},
			err:  nil,
		},
		{
			name: "negativ check",
			metr: common.Metrics{ID: "testC2", Delta: &testC2V, MType: "counter"},
			err:  nil,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBasicStorage := mock.NewMockMetricsStorage(ctrl)

			if test.metr.MType == "gauge" {
				mockBasicStorage.EXPECT().SetGaugeValue(test.metr.ID, common.TypeGauge(*test.metr.Value)).Return(nil).Times(1)
			}
			if test.metr.MType == "counter" {
				mockBasicStorage.EXPECT().SetCounterValue(test.metr.ID, common.TypeCounter(*test.metr.Delta)).Return(nil).Times(1)
			}

			mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}

			err := mServ.updateValue(&test.metr)
			if err != test.err {
				t.Errorf("Problem error. actual:%s;expected:%s", err, test.err)
			}
		})
	}
}

func TestUpdateValues(t *testing.T) {
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockMetricsStorage(ctrl)
	mockBasicStorage.EXPECT().AddMetrics(testMetrics).Return(nil).Times(1)
	mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}

	err := mServ.updateValues(testMetrics)

	require.NoError(t, err)
}

func TestUpdateJSONValue(t *testing.T) {
	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []float64{6.7, 7.89}
	metricsCounterName := []string{"testCounter1", "testCounter2"}
	metricsCounterval := []int64{56, 12}

	testData := []struct {
		name   string
		data   common.Metrics
		status int
	}{
		{
			name:   "check 1",
			data:   common.Metrics{ID: metricsGaugeName[0], MType: "gauge", Value: &metricsGaugeVal[0]},
			status: http.StatusOK,
		},
		{
			name:   "check 2",
			data:   common.Metrics{ID: metricsGaugeName[1], MType: "gauge", Value: &metricsGaugeVal[1]},
			status: http.StatusOK,
		},
		{
			name:   "check 3",
			data:   common.Metrics{ID: metricsCounterName[0], MType: "counter", Delta: &metricsCounterval[0]},
			status: http.StatusOK,
		},
		{
			name:   "check 4",
			data:   common.Metrics{ID: metricsCounterName[1], MType: "counter", Delta: &metricsCounterval[1]},
			status: http.StatusOK,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBasicStorage := mock.NewMockMetricsStorage(ctrl)

			if test.data.MType == "gauge" {
				mockBasicStorage.EXPECT().SetGaugeValue(test.data.ID, common.TypeGauge(*test.data.Value)).Return(nil).Times(1)
			}
			if test.data.MType == "counter" {
				mockBasicStorage.EXPECT().SetCounterValue(test.data.ID, common.TypeCounter(*test.data.Delta)).Return(nil).Times(1)
			}
			d, _ := json.Marshal(test.data)
			myReader := strings.NewReader(string(d))
			request := httptest.NewRequest(http.MethodGet, "/update/", myReader)
			w := httptest.NewRecorder()
			mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}
			expectdStatus := http.StatusOK

			mServ.UpdateJSONValue(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, expectdStatus, res.StatusCode)
		})
	}
}

func TestUpdatesMetrics(t *testing.T) {
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBasicStorage := mock.NewMockMetricsStorage(ctrl)
	mockBasicStorage.EXPECT().AddMetrics(testMetrics).Return(nil).Times(1)

	d, _ := json.Marshal(testMetrics)
	myReader := strings.NewReader(string(d))
	request := httptest.NewRequest(http.MethodGet, "/updates/", myReader)
	w := httptest.NewRecorder()
	mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}
	expectdStatus := http.StatusOK

	mServ.UpdatesMetrics(w, request)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, expectdStatus, res.StatusCode)
}

func TestGetValue(t *testing.T) {
	testG1V := 54.7
	testC2V := int64(5)
	testData := []struct {
		name   string
		metr   common.Metrics
		status int
	}{
		{
			name:   "positiv check",
			metr:   common.Metrics{ID: "testG1", Value: &testG1V, MType: "gauge"},
			status: http.StatusOK,
		},
		{
			name:   "negativ check",
			metr:   common.Metrics{ID: "testC2", Delta: &testC2V, MType: "counter"},
			status: http.StatusOK,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBasicStorage := mock.NewMockMetricsStorage(ctrl)

			if test.metr.MType == "gauge" {
				mockBasicStorage.EXPECT().GetGaugeValue(test.metr.ID).Return(common.TypeGauge(3.0), nil).Times(1)
			}
			if test.metr.MType == "counter" {
				mockBasicStorage.EXPECT().GetCounterValue(test.metr.ID).Return(common.TypeCounter(4), nil).Times(1)
			}

			mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}

			st := mServ.getValue(&test.metr)
			assert.Equal(t, test.status, st)
		})
	}
}

func TestGetJSONValue(t *testing.T) {
	metricsGaugeName := []string{"testGauge1", "testGauge2"}
	metricsGaugeVal := []float64{6.7, 7.89}

	testData := []struct {
		name     string
		data     common.Metrics
		status   int
		expected string
	}{
		{
			name:     "check 1",
			data:     common.Metrics{ID: metricsGaugeName[0], MType: "gauge"},
			status:   http.StatusOK,
			expected: `{"value":6.7,"id":"testGauge1","type":"gauge"}`,
		},
		{
			name:     "check 2",
			data:     common.Metrics{ID: metricsGaugeName[1], MType: "gauge"},
			status:   http.StatusOK,
			expected: `{"value":7.89,"id":"testGauge2","type":"gauge"}`,
		},
	}

	for i, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBasicStorage := mock.NewMockMetricsStorage(ctrl)

			mockBasicStorage.EXPECT().GetGaugeValue(test.data.ID).Return(common.TypeGauge(metricsGaugeVal[i]), nil).Times(1)

			d, _ := json.Marshal(test.data)
			myReader := strings.NewReader(string(d))
			request := httptest.NewRequest(http.MethodGet, "/value/", myReader)
			w := httptest.NewRecorder()
			mServ := MetricServer{logger: &common.FakeLogger{}, memStorage: mockBasicStorage, signKey: signKey}
			expectdStatus := http.StatusOK

			mServ.GetJSONValue(w, request)

			res := w.Result()
			dataT := make([]byte, 300)
			n, _ := res.Body.Read(dataT)
			res.Body.Close()
			assert.Equal(t, expectdStatus, res.StatusCode)
			assert.Equal(t, test.expected, string(dataT[:n]))
		})
	}
}
