package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/Alexandrfield/Metrics/internal/requestHandler"
	"github.com/Alexandrfield/Metrics/internal/server"
	"github.com/Alexandrfield/Metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUpdateValue(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "/update/counter/test/1",
			want: want{
				code:        200,
				contentType: "plain/text",
			},
		},
		{
			name:    "positive test #2",
			request: "/update/gauge/test/4.4",
			want: want{
				code:        200,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #2",
			request: "/update/counter/5",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #3",
			request: "/update/gauge/5",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #4",
			request: "/update/counter/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #5",
			request: "/update/gauge/",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #6",
			request: "/update/unknown/testCounter/100",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #7",
			request: "/update/counter/testCounter/none",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #8",
			request: "/update/gauge/testCounter/none",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "plain/text",
			},
		},
	}
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := storage.Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	store := storage.CreateMemStorage(storageConfig, logger, done)
	metricRep := server.CreateMetricRepository(store, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, http.NoBody)
			w := httptest.NewRecorder()

			servHandler.UpdateValue(w, request)
			result := w.Result()
			_ = result.Body.Close()
			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}
func TestDefaultAnswer(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/update/unknown/test/1", http.NoBody)
	w := httptest.NewRecorder()

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()

	done := make(chan struct{})
	storageConfig := storage.Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	store := storage.CreateMemStorage(storageConfig, logger, done)

	metricRep := server.CreateMetricRepository(store, logger)
	servHandler := handler.CreateHandlerRepository(&metricRep, logger)

	servHandler.DefaultAnswer(w, request)
	result := w.Result()
	assert.Equal(t, http.StatusNotImplemented, result.StatusCode)
	_ = result.Body.Close()
}

func TestParserURL(t *testing.T) {
	tests := []struct {
		name      string
		request   string
		resStatus int
	}{
		{
			name:      "positive test #1",
			request:   "/update/counter/test/1",
			resStatus: http.StatusOK,
		},
		{
			name:      "negative test #1",
			request:   "/update/counter/1",
			resStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, http.NoBody)
			w := httptest.NewRecorder()

			zapLogger, err := zap.NewDevelopment()
			if err != nil {
				t.Errorf("Can not initializate zap logger. err:%v", err)
				return
			}
			defer func() { _ = zapLogger.Sync() }()
			logger := zapLogger.Sugar()

			done := make(chan struct{})
			storageConfig := storage.Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
			store := storage.CreateMemStorage(storageConfig, logger, done)
			metricRep := server.CreateMetricRepository(store, logger)
			servHandler := handler.CreateHandlerRepository(&metricRep, logger)

			servHandler.UpdateValue(w, request)
			result := w.Result()
			assert.Equal(t, tt.resStatus, result.StatusCode)
			_ = result.Body.Close()
		})
	}
}
