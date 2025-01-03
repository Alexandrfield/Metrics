package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

var errMetricNotExistIssue = errors.New("metric with this name or type is does't exist")

type TypeGauge float64
type TypeCounter int64

type MemStorage struct {
	gaugeData   map[string]TypeGauge
	counterData map[string]TypeCounter
	logger      common.Loger
	Config      Config
}

func CreateMemStorage(config Config, logger common.Loger, done chan struct{}) *MemStorage {
	memStorage := MemStorage{gaugeData: make(map[string]TypeGauge),
		counterData: make(map[string]TypeCounter), logger: logger, Config: config}
	logger.Debugf("config.Restore %s", config.Restore)
	if config.Restore {
		file, err := os.OpenFile(memStorage.Config.FileStoregePath, os.O_RDONLY, 0o600)
		if err == nil {
			memStorage.logger.Debugf("file was open")
			defer func() {
				_ = file.Close()
			}()
			memStorage.LoadMemStorage(file)
		} else {
			logger.Debugf("can not restore file. File is not exist. err:%w", err)
		}
	}
	if config.StoreIntervalSecond != 0 {
		go storageSaver(&memStorage, config.FileStoregePath, config.StoreIntervalSecond, done)
	}
	return &memStorage
}
func storageSaver(memStorage *MemStorage, filepath string, saveIntervalSecond int, done chan struct{}) {
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
func (st *MemStorage) saveMemStorageInFile(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		st.logger.Debugf("Issue with open %s %w", filename, err)
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
func (st *MemStorage) saveMemStorage(stream io.Writer) {
	for key, val := range st.gaugeData {
		temp := float64(val)
		metric := common.Metrics{ID: key, MType: "gauge", Value: &temp}
		_, _ = stream.Write([]byte(createStringMetric("gauge", key, metric.GetValueMetric())))
	}
	for key, val := range st.counterData {
		temp := int64(val)
		metric := common.Metrics{ID: key, MType: "counter", Delta: &temp}
		_, _ = stream.Write([]byte(createStringMetric("counter", key, metric.GetValueMetric())))
	}
}
func (st *MemStorage) LoadMemStorage(stream io.Reader) {
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
			st.logger.Debugf("metric->%v; %s,%s,%s;", metric, res[0], res[1], res[2])
			if err != nil {
				st.logger.Debugf("metric.SaveMetric err:%v;", err)
			}
			switch metric.MType {
			case "gauge":
				_ = st.AddGauge(metric.ID, TypeGauge(*metric.Value))
			case "counter":
				_ = st.AddCounter(metric.ID, TypeCounter(*metric.Delta))
			}
		}
	}
}

func (st *MemStorage) AddGauge(name string, value TypeGauge) error {
	st.gaugeData[name] = value
	return nil
}
func (st *MemStorage) GetGauge(name string) (TypeGauge, error) {
	val, ok := st.gaugeData[name]
	if !ok {
		return TypeGauge(0), fmt.Errorf("can't find Gauge metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemStorage) AddCounter(name string, value TypeCounter) error {
	val, ok := st.counterData[name]
	if !ok {
		val = 0
	}
	st.counterData[name] = val + value
	return nil
}
func (st *MemStorage) GetCounter(name string) (TypeCounter, error) {
	val, ok := st.counterData[name]
	if !ok {
		return TypeCounter(0), fmt.Errorf("can't find Counter metric with name:%s;err:%w", name, errMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemStorage) GetAllMetricName() ([]string, []string) {
	allGaugeKeys := make([]string, 0)
	for key := range st.gaugeData {
		allGaugeKeys = append(allGaugeKeys, key)
	}
	allCounterKeys := make([]string, 0)
	for key := range st.counterData {
		allCounterKeys = append(allCounterKeys, key)
	}
	return allGaugeKeys, allCounterKeys
}
