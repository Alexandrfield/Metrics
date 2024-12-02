package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGaugePositiv(t *testing.T) {
	memStorage := CreateMemStorage()

	tests := []struct {
		name   string
		value  string
		expect TypeGauge
	}{
		{
			name:   "test1",
			value:  "24",
			expect: 24,
		},
		{
			name:   "test2",
			value:  "-24",
			expect: -24,
		},
		{
			name:   "test3",
			value:  "24.5",
			expect: 24.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage.AddGauge(tt.name, tt.value)
			res, ok := memStorage.GetGauge(tt.name)
			assert.Equal(t, ok, true)
			assert.Equal(t, tt.expect, res)
		})
	}

}
func TestAddGaugeNegativ(t *testing.T) {
	memStorage := CreateMemStorage()

	memStorage.AddGauge("test1", "23")
	_, ok := memStorage.GetGauge("test2")
	assert.Equal(t, ok, false)
}

func TestAddCounterPositiv(t *testing.T) {
	memStorage := CreateMemStorage()

	tests := []struct {
		name   string
		value  string
		expect TypeCounter
	}{
		{
			name:   "test1",
			value:  "24",
			expect: 24,
		},
		{
			name:   "test2",
			value:  "-24",
			expect: -24,
		},
		{
			name:   "test3",
			value:  "0",
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage.AddCounter(tt.name, tt.value)
			res, ok := memStorage.GetCounter(tt.name)
			assert.Equal(t, ok, true)
			assert.Equal(t, tt.expect, res)
		})
	}

}
func TestAddCounterNegativ(t *testing.T) {
	memStorage := CreateMemStorage()

	memStorage.AddCounter("test1", "23")
	_, ok := memStorage.GetCounter("test2")
	assert.Equal(t, ok, false)
}
