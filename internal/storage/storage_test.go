package storage

import (
	"errors"
	"log"
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
	memStorage := CreateMemStorage()

	err := memStorage.AddGauge("test1", "23")
	if err != nil {
		t.Errorf("Error for test; err:%s\n", err)
		return
	}
	_, err = memStorage.GetGauge("test2")
	check := errors.Is(err, errMetricNotExistIssue)
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
			value:  "42",
			expect: "42",
		},
		{
			name:   "test2",
			value:  "-77",
			expect: "-77",
		},
		{
			name:   "test3",
			value:  "0",
			expect: "0",
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
	memStorage := CreateMemStorage()

	err := memStorage.AddCounter("test1", "23")
	if err != nil {
		t.Errorf("Error for test; err:%s\n", err)
		return
	}
	_, err = memStorage.GetCounter("test2")

	check := errors.Is(err, errMetricNotExistIssue)
	assert.Equal(t, check, true)
}

func TestGetAllMetricName(t *testing.T) {
	memStorage := CreateMemStorage()

	testsGauge := []struct {
		name  string
		value string
	}{
		{
			name:  "testsGauge1",
			value: "14",
		},
		{
			name:  "testsGauge2",
			value: "-14",
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
