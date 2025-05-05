package storage

import (
	"errors"
	"log"
	"testing"

	"github.com/Alexandrfield/Metrics/internal/common"
	file_store "github.com/Alexandrfield/Metrics/internal/storage/file"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddGaugePositiv(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	memStorage := CreateMemStorage(storageConfig, logger, done)

	tests := []struct {
		name   string
		value  common.TypeGauge
		expect common.TypeGauge
	}{
		{
			name:   "test1",
			value:  24,
			expect: 24,
		},
		{
			name:   "test2",
			value:  -24,
			expect: -24,
		},
		{
			name:   "test3",
			value:  24.5,
			expect: 24.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := memStorage.AddGauge(tt.name, tt.value)
			if err != nil {
				log.Printf("Error for test; err:%v\n", err)
			}
			res, err := memStorage.GetGauge(tt.name)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.expect, res)
		})
	}
}
func TestAddGaugeNegativ(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	memStorage := CreateMemStorage(storageConfig, logger, done)

	err = memStorage.AddGauge("test1", 23)
	if err != nil {
		t.Errorf("Error for test; err:%s\n", err)
		return
	}
	_, err = memStorage.GetGauge("test2")
	check := errors.Is(err, file_store.ErrMetricNotExistIssue)
	assert.Equal(t, check, true)
}

func TestAddCounterPositiv(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	memStorage := CreateMemStorage(storageConfig, logger, done)

	tests := []struct {
		name   string
		value  common.TypeCounter
		expect common.TypeCounter
	}{
		{
			name:   "test1",
			value:  42,
			expect: 42,
		},
		{
			name:   "test2",
			value:  -77,
			expect: -77,
		},
		{
			name:   "test3",
			value:  0,
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := memStorage.AddCounter(tt.name, tt.value)
			if err != nil {
				t.Errorf("Error for test; err:%s\n", err)
				return
			}
			res, err := memStorage.GetCounter(tt.name)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.expect, res)
		})
	}
}
func TestAddCounterNegativ(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	memStorage := CreateMemStorage(storageConfig, logger, done)

	err = memStorage.AddCounter("test1", 23)
	if err != nil {
		t.Errorf("Error for test; err:%s\n", err)
		return
	}
	_, err = memStorage.GetCounter("test2")

	check := errors.Is(err, file_store.ErrMetricNotExistIssue)
	assert.Equal(t, check, true)
}

func TestGetAllMetricName(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Errorf("Can not initializate zap logger. err:%v", err)
		return
	}
	defer func() { _ = zapLogger.Sync() }()
	logger := zapLogger.Sugar()
	done := make(chan struct{})
	storageConfig := Config{FileStoregePath: "test.log", StoreIntervalSecond: 0, Restore: false}
	memStorage := CreateMemStorage(storageConfig, logger, done)

	testsGauge := []struct {
		name  string
		value common.TypeGauge
	}{
		{
			name:  "testsGauge1",
			value: 14,
		},
		{
			name:  "testsGauge2",
			value: -14,
		},
		{
			name:  "testsGauge3",
			value: 0,
		},
	}
	testsCounter := []struct {
		name  string
		value common.TypeCounter
	}{
		{
			name:  "testCounter1",
			value: 24,
		},
		{
			name:  "testCounter2",
			value: -24,
		},
		{
			name:  "testCounter3",
			value: 0,
		},
	}
	expected := make([]string, 0)
	actual := make([]string, 0)
	for _, tt := range testsGauge {
		err := memStorage.AddGauge(tt.name, tt.value)
		if err != nil {
			t.Errorf("Error for test; err:%s\n", err)
			return
		}
		expected = append(expected, tt.name)
	}
	for _, tt := range testsCounter {
		err := memStorage.AddCounter(tt.name, tt.value)
		if err != nil {
			t.Errorf("Error for test; err:%s\n", err)
			return
		}
		expected = append(expected, tt.name)
	}

	keysCounter, keysGauge := memStorage.GetAllMetricName()
	actual = append(actual, keysCounter...)
	actual = append(actual, keysGauge...)
	assert.ElementsMatch(t, actual, expected)
}

func TestCreateMemStorage(t *testing.T) {
	storageConfig := Config{DatabaseDsn: "testerr", StoreIntervalSecond: 0, Restore: false}
	done := make(chan struct{})
	memStorage := CreateMemStorage(storageConfig, &common.FakeLogger{}, done)
	if memStorage == nil {
		t.Errorf("problem with create db")
	}
}
func TestCreateMemStorageSecond(t *testing.T) {
	storageConfig := Config{FileStoregePath: "", StoreIntervalSecond: 3, Restore: true}
	done := make(chan struct{})
	memStorage := CreateMemStorage(storageConfig, &common.FakeLogger{}, done)
	if memStorage == nil {
		t.Errorf("problem with create db file")
	}
	memStorage.Close()
}
