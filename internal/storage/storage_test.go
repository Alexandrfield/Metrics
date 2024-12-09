package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGaugePositiv(t *testing.T) {
	memStorage := CreateMemStorage()

	tests := []struct {
		name   string
		value  string
		expect string
	}{
		{
			name:   "test1",
			value:  "24",
			expect: "24",
		},
		{
			name:   "test2",
			value:  "-24",
			expect: "-24",
		},
		{
			name:   "test3",
			value:  "24.5",
			expect: "24.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage.AddGauge(tt.name, tt.value)
			res, err := memStorage.GetGauge(tt.name)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.expect, res)
		})
	}

}
func TestAddGaugeNegativ(t *testing.T) {
	memStorage := CreateMemStorage()

	memStorage.AddGauge("test1", "23")
	_, err := memStorage.GetGauge("test2")
	check := errors.Is(err, ErrMetricNotExistIssue)
	assert.Equal(t, check, true)
}

func TestAddCounterPositiv(t *testing.T) {
	memStorage := CreateMemStorage()

	tests := []struct {
		name   string
		value  string
		expect string
	}{
		{
			name:   "test1",
			value:  "24",
			expect: "24",
		},
		{
			name:   "test2",
			value:  "-24",
			expect: "-24",
		},
		{
			name:   "test3",
			value:  "0",
			expect: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage.AddCounter(tt.name, tt.value)
			res, err := memStorage.GetCounter(tt.name)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.expect, res)
		})
	}
}
func TestAddCounterNegativ(t *testing.T) {
	memStorage := CreateMemStorage()

	memStorage.AddCounter("test1", "23")
	_, err := memStorage.GetCounter("test2")

	check := errors.Is(err, ErrMetricNotExistIssue)
	assert.Equal(t, check, true)
}

func TesGetAllMetricName(t *testing.T) {
	memStorage := CreateMemStorage()

	testsGauge := []struct {
		name  string
		value string
	}{
		{
			name:  "testsGauge1",
			value: "24",
		},
		{
			name:  "testsGauge2",
			value: "-24",
		},
		{
			name:  "testsGauge3",
			value: "0",
		},
	}
	testsCounter := []struct {
		name  string
		value string
	}{
		{
			name:  "testCounter1",
			value: "24",
		},
		{
			name:  "testCounter2",
			value: "-24",
		},
		{
			name:  "testCounter3",
			value: "0",
		},
	}
	var expected []string
	var actual []string
	for _, tt := range testsGauge {
		memStorage.AddGauge(tt.name, tt.value)
		expected = append(expected, tt.name)
	}
	for _, tt := range testsCounter {
		memStorage.AddCounter(tt.name, tt.value)
		expected = append(expected, tt.name)
	}

	keysCounter, keysGauge := memStorage.GetAllMetricName()
	actual = append(actual, keysCounter...)
	actual = append(actual, keysGauge...)
	assert.ElementsMatch(t, actual, expected)
}
