package filestorage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Alexandrfield/Metrics/internal/common"
)

var ErrMetricNotExistIssue = errors.New("metric with this name or type is does't exist")
var ErrObjectHasbeenClosed = errors.New("metric storage MemFileStorage has been already closed")

const typecounter = "counter"
const typegauge = "gauge"

type MemFileStorage struct {
	GaugeData   map[string]common.TypeGauge
	CounterData map[string]common.TypeCounter
	Logger      common.Loger
	filepath    string
	isCreated   bool
	lock        sync.Mutex
}

func StorageSaver(memStorage *MemFileStorage, saveIntervalSecond int, done chan struct{}) {
	if saveIntervalSecond != 0 {
		tickerSaveInterval := time.NewTicker(time.Duration(saveIntervalSecond) * time.Second)
		for {
			select {
			case <-done:
				memStorage.close()
				return
			case <-tickerSaveInterval.C:
				memStorage.saveMemStorageInFile()
			}
		}
	} else {
		<-done
		memStorage.close()
	}
}
func (st *MemFileStorage) Close() {
	// action will be done in StorageSaver
}
func (st *MemFileStorage) close() {
	if st.isCreated {
		st.isCreated = false
		st.lock.Lock()
		st.saveMemStorageInFile()
	}
}
func (st *MemFileStorage) saveMemStorageInFile() {
	file, err := os.OpenFile(st.filepath, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		st.Logger.Debugf("Issue with open %s %w", st.filepath, err)
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
func NewMemFileStorage(filepath string, logger common.Loger) *MemFileStorage {
	memStorage := MemFileStorage{GaugeData: make(map[string]common.TypeGauge),
		CounterData: make(map[string]common.TypeCounter), isCreated: true, filepath: filepath, Logger: logger}
	return &memStorage
}

func (st *MemFileStorage) saveMemStorage(stream io.Writer) {
	if !st.isCreated {
		return
	}
	st.lock.Lock()
	defer st.lock.Unlock()
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
	if !st.isCreated {
		return
	}
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
	if !st.isCreated {
		return ErrObjectHasbeenClosed
	}
	st.lock.Lock()
	defer st.lock.Unlock()
	st.GaugeData[name] = value
	return nil
}
func (st *MemFileStorage) GetGauge(name string) (common.TypeGauge, error) {
	if !st.isCreated {
		return common.TypeGauge(0), ErrObjectHasbeenClosed
	}
	st.lock.Lock()
	defer st.lock.Unlock()
	val, ok := st.GaugeData[name]
	if !ok {
		return common.TypeGauge(0), fmt.Errorf("can't find Gauge metric with name:%s;err:%w", name, ErrMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemFileStorage) AddCounter(name string, value common.TypeCounter) error {
	if !st.isCreated {
		return ErrObjectHasbeenClosed
	}
	st.lock.Lock()
	defer st.lock.Unlock()
	val, ok := st.CounterData[name]
	if !ok {
		val = 0
	}
	st.CounterData[name] = val + value
	return nil
}
func (st *MemFileStorage) GetCounter(name string) (common.TypeCounter, error) {
	if !st.isCreated {
		return common.TypeCounter(0), ErrObjectHasbeenClosed
	}
	st.lock.Lock()
	defer st.lock.Unlock()
	val, ok := st.CounterData[name]
	if !ok {
		return common.TypeCounter(0), fmt.Errorf("can't find Counter metric with name:%s;err:%w",
			name, ErrMetricNotExistIssue)
	}
	return val, nil
}
func (st *MemFileStorage) GetAllMetricName() ([]string, []string) {
	if !st.isCreated {
		return []string{}, []string{}
	}
	st.lock.Lock()
	defer st.lock.Unlock()
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
	if !st.isCreated {
		return ErrObjectHasbeenClosed
	}
	for _, metric := range metrics {
		switch metric.MType {
		case typegauge:
			_ = st.AddGauge(metric.ID, common.TypeGauge(*metric.Value))
		case typecounter:
			_ = st.AddCounter(metric.ID, common.TypeCounter(*metric.Delta))
		default:
			return fmt.Errorf("AddMetrics. unknown type:%s;", metric.MType)
		}
	}
	return nil
}
