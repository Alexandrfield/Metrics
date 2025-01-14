package filestorage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

var ErrMetricNotExistIssue = errors.New("metric with this name or type is does't exist")

const typecounter = "counter"
const typegauge = "gauge"

type MemFileStorage struct {
	GaugeData   map[string]common.TypeGauge
	CounterData map[string]common.TypeCounter
	Logger      common.Loger
}

func StorageSaver(memStorage *MemFileStorage, filepath string, saveIntervalSecond int, done chan struct{}) {
	tickerSaveInterval := time.NewTicker(time.Duration(saveIntervalSecond) * time.Second)
	for {
		select {
		case <-done:
			memStorage.saveMemStorageInFile(filepath)
			return
		case <-tickerSaveInterval.C:
			memStorage.saveMemStorageInFile(filepath)
		}
	}
}
func (st *MemFileStorage) saveMemStorageInFile(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		st.Logger.Debugf("Issue with open %s %w", filename, err)
		return
	}
	defer func() {
		_ = file.Close()
	}()
	st.saveMemStorage(file)
}

func createStringMetric(mtype string, name string, value string) string {
	return fmt.Sprintf("%s;%s;%s\n", mtype, name, value)
}
func (st *MemFileStorage) saveMemStorage(stream io.Writer) {
	for key, val := range st.GaugeData {
		temp := float64(val)
		metric := common.Metrics{ID: key, MType: typegauge, Value: &temp}
		_, _ = stream.Write([]byte(createStringMetric(typegauge, key, metric.GetValueMetric())))
	}
	for key, val := range st.CounterData {
		temp := int64(val)
		metric := common.Metrics{ID: key, MType: typecounter, Delta: &temp}
		_, _ = stream.Write([]byte(createStringMetric(typecounter, key, metric.GetValueMetric())))
	}
}
func (st *MemFileStorage) LoadMemStorage(stream io.Reader) {
	data := make([]byte, 1000)
	for {
		n, err := stream.Read(data)
		if errors.Is(err, io.EOF) {
			break
		}
		rawData := string(data[:n])
		listMetrics := strings.Split(rawData, "\n")
		for _, tmp := range listMetrics {
			res := strings.Split(tmp, ";")
			if len(res) < 3 {
				continue
			}
			var metric common.Metrics
			err = metric.SaveMetric(res[0], res[1], res[2])
			st.Logger.Debugf("metric->%v; %s,%s,%s;", metric, res[0], res[1], res[2])
			if err != nil {
				st.Logger.Debugf("metric.SaveMetric err:%v;", err)
			}
			switch metric.MType {
			case typegauge:
				_ = st.AddGauge(metric.ID, common.TypeGauge(*metric.Value))
			case typecounter:
				_ = st.AddCounter(metric.ID, common.TypeCounter(*metric.Delta))
			}
		}
	}
}

func (st *MemFileStorage) AddGauge(name string, value common.TypeGauge) error {
	st.GaugeData[name] = value
	return nil
}
func (st *MemFileStorage) GetGauge(name string) (common.TypeGauge, error) {
	val, ok := st.GaugeData[name]
	if !ok {
		return common.TypeGauge(0), fmt.Errorf("can't find Gauge metric with name:%s;err:%w", name, ErrMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemFileStorage) AddCounter(name string, value common.TypeCounter) error {
	val, ok := st.CounterData[name]
	if !ok {
		val = 0
	}
	st.CounterData[name] = val + value
	return nil
}
func (st *MemFileStorage) GetCounter(name string) (common.TypeCounter, error) {
	val, ok := st.CounterData[name]
	if !ok {
		return common.TypeCounter(0), fmt.Errorf("can't find Counter metric with name:%s;err:%w",
			name, ErrMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemFileStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	for key := range st.GaugeData {
		allGaugeKeys = append(allGaugeKeys, key)
	}
	allCounterKeys := make([]string, 0)
	for key := range st.CounterData {
		allCounterKeys = append(allCounterKeys, key)
	}
	return allGaugeKeys, allCounterKeys
}

func (st *MemFileStorage) PingDatabase() bool {
	return false
}

func (st *MemFileStorage) AddMetrics(metrics []common.Metrics) error {
	for _, metric := range metrics {
		var err error
		switch metric.MType {
		case typegauge:
			err = st.AddGauge(metric.ID, common.TypeGauge(*metric.Value))
		case typecounter:
			err = st.AddCounter(metric.ID, common.TypeCounter(*metric.Delta))
		default:
			return fmt.Errorf("AddMetrics. unknown type:%s;", metric.MType)
		}
		if err != nil {
			st.Logger.Debugf("Problem with add Metrics. err:%s", err)
		}
	}
	return nil
}
