package requesthandler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Alexandrfield/Metrics/internal/common"
	"github.com/Alexandrfield/Metrics/internal/storage"
)

type MetricsStorage interface {
	SetCounterValue(metricName string, metricValue storage.TypeCounter) error
	SetGaugeValue(metricName string, metricValue storage.TypeGauge) error
	GetCounterValue(metricName string) (storage.TypeCounter, error)
	GetGaugeValue(metricName string) (storage.TypeGauge, error)
	GetAllValue() ([]string, error)
}

type MetricServer struct {
	logger     common.Loger
	memStorage MetricsStorage
}

func CreateHandlerRepository(stor MetricsStorage) *MetricServer {
	return &MetricServer{memStorage: stor}
}

func (rep *MetricServer) DefaultAnswer(res http.ResponseWriter, req *http.Request) {
	rep.logger.Debugf("defaultAnswer. req:%v;res.WriteHeader::%d\n", req, http.StatusNotImplemented)
	res.WriteHeader(http.StatusNotImplemented)
}

func (rep *MetricServer) UpdateValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusOK
	var metric common.Metrics
	var err error
	if err = json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	rep.logger.Debugf("setValue type:%s; name%s; value:%d; delta:%d; err:%s\n", metric.MType, metric.ID, metric.Value, metric.Delta, err)
	switch metric.MType {
	case "gauge":
		err = rep.memStorage.SetGaugeValue(metric.ID, storage.TypeGauge(*metric.Value))
	case "counter":
		err = rep.memStorage.SetCounterValue(metric.ID, storage.TypeCounter(*metric.Delta))
	default:
		statusH = http.StatusBadRequest
		rep.logger.Warnf("unknown type:%s", metric.MType)
	}
	if err != nil {
		rep.logger.Warnf("unknown type:%s", metric.MType)
	}
	res.WriteHeader(statusH)
}
func (rep *MetricServer) GetValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusOK
	var metric common.Metrics
	if err := json.NewDecoder(req.Body).Decode(&metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	rep.logger.Debugf("getValue type:%s; name%s; value:%d; delta:%d;", metric.MType, metric.ID, metric.Value, metric.Delta)

	//	var err error
	switch metric.MType {
	case "gauge":
		val, err := rep.memStorage.GetGaugeValue(metric.ID)
		temp := float64(val)
		metric.Value = &temp
		if err != nil {
			statusH = http.StatusNotFound
		}
	case "counter":
		val, err := rep.memStorage.GetCounterValue(metric.ID)
		temp := int64(val)
		metric.Delta = &temp
		if err != nil {
			statusH = http.StatusNotFound
		}
	default:
		statusH = http.StatusBadRequest
		rep.logger.Warnf("unknown type:%s", metric.MType)
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		rep.logger.Warnf("problem with unmarshal:%w", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusH)
	res.Write(resp)
}

func (rep *MetricServer) GetAllData(res http.ResponseWriter, req *http.Request) {
	allValues, err := rep.memStorage.GetAllValue()
	if err != nil {
		rep.logger.Debugf("issue for GetAllData. err:%s\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	allValuesTemplate := `
<html> 
   <head> 
   </head> 
   <body> 
	all metrics: 
		{{ range .}}{{.}}, {{ end }}
   </body> 
</html>
`
	readyTemplate, err := template.New("templ").Parse(allValuesTemplate)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	err = readyTemplate.Execute(res, allValues)
	if err != nil {
		rep.logger.Debugf("issue for readyTemplate.Execute(res, allValues). err:%s\n", err)
	}
}
