package common

import (
	"testing"
)

func TestSaveMetric(t *testing.T) {
	var temp1 int64 = 45
	var temp2 int64 = 3
	temp3 := 22.5
	temp4 := 13.03
	testData := []Metrics{
		{ID: "test1", MType: "counter", Delta: &temp1},
		{ID: "test2", MType: "counter", Delta: &temp2},
		{ID: "test3", MType: "gauge", Value: &temp3},
		{ID: "test4", MType: "gauge", Value: &temp4},
	}

	for _, v := range testData {
		temp := Metrics{}
		err := temp.SaveMetric(v.MType, v.ID, v.GetValueMetric())
		if err != nil {
			t.Errorf("error Metric name:%s Delta. actual:%s, expected err:nil", temp.ID, err)
			return
		}
		if temp.Delta != nil {
			if *temp.Delta != *v.Delta {
				t.Errorf("error Metric name:%s Delta. actual:%d, expected:%d", temp.ID, *temp.Delta, *v.Delta)
				return
			}
		}
		if temp.Value != nil {
			if *temp.Value != *v.Value {
				t.Errorf("error Metric name:%s Metric Value. actual:%f, expected:%f", temp.ID, *temp.Value, *v.Value)
				return
			}
		}
		if temp.MType != v.MType {
			t.Errorf("error Metric name:%s Metric MType. actual:%s, expected:%s", temp.ID, temp.MType, v.MType)
			return
		}
		if temp.ID != v.ID {
			t.Errorf("error Metric name:%s Metric Value. actual:%s, expected:%s", temp.ID, temp.ID, v.ID)
			return
		}
	}
}

func TestGetValueMetric(t *testing.T) {
	var temp1 int64 = 45
	var temp2 int64 = 3
	temp3 := 22.5
	temp4 := 13.03
	expectedStr := []string{"45", "3", "22.5", "13.03"}
	testData := []Metrics{
		{ID: "test1", MType: "counter", Delta: &temp1},
		{ID: "test2", MType: "counter", Delta: &temp2},
		{ID: "test3", MType: "gauge", Value: &temp3},
		{ID: "test4", MType: "gauge", Value: &temp4},
	}
	for i, v := range testData {
		temp := v.GetValueMetric()
		if expectedStr[i] != v.GetValueMetric() {
			t.Errorf("error Metric name:%s Metric string data. actual:%s, expected:%s", v.ID, temp, expectedStr[i])
			return
		}
	}
}
