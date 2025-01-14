package requesthandler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Alexandrfield/Metrics/internal/common"
)

type MetricsStorage interface {
	SetCounterValue(metricName string, metricValue common.TypeCounter) error
	SetGaugeValue(metricName string, metricValue common.TypeGauge) error
	GetCounterValue(metricName string) (common.TypeCounter, error)
	GetGaugeValue(metricName string) (common.TypeGauge, error)
	GetAllValue() ([]string, error)
	PingDatabase() bool
	AddMetrics(metrics []common.Metrics) error
}

type MetricServer struct {
	logger     common.Loger
	memStorage MetricsStorage
}

func parseURL(req *http.Request, logger common.Loger) (common.Metrics, int) {
	var metric common.Metrics
	url := strings.Split(req.URL.String(), "/")
	// expected format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>,
	// Content-Type: text/plain
	if url[1] == "update" {
		if len(url) != 5 {
			return metric, http.StatusNotFound
		}
		err := metric.SaveMetric(url[2], url[3], url[4])
		if err != nil {
			logger.Debugf("issue with parse metric (command update): %w", err)
			return metric, http.StatusBadRequest
		}
		return metric, http.StatusOK
	}
	// expected format http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	if url[1] == "value" {
		if len(url) != 4 {
			return metric, http.StatusNotFound
		}
		metric.MType = url[2]
		metric.ID = url[3]
		return metric, http.StatusOK
	}
	return metric, http.StatusNotFound
}

func CreateHandlerRepository(stor MetricsStorage, logger common.Loger) *MetricServer {
	return &MetricServer{memStorage: stor, logger: logger}
}

func (rep *MetricServer) DefaultAnswer(res http.ResponseWriter, req *http.Request) {
	rep.logger.Debugf("defaultAnswer. req:%v;res.WriteHeader::%d\n", req, http.StatusNotImplemented)
	res.WriteHeader(http.StatusNotImplemented)
}

func (rep *MetricServer) Ping(res http.ResponseWriter, req *http.Request) {
	rep.logger.Debugf("pingDatabase. req:%v;", req)
	if rep.memStorage.PingDatabase() {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (rep *MetricServer) updateValue(metric *common.Metrics) error {
	rep.logger.Debugf("setValue type:%s; name: %s; value:%d; delta:%d;",
		metric.MType, metric.ID, metric.Value, metric.Delta)
	var err error
	switch metric.MType {
	case "gauge":
		err = rep.memStorage.SetGaugeValue(metric.ID, common.TypeGauge(*metric.Value))
	case "counter":
		err = rep.memStorage.SetCounterValue(metric.ID, common.TypeCounter(*metric.Delta))
	default:
		return fmt.Errorf("unknown type:%s;", metric.MType)
	}
	if err != nil {
		return fmt.Errorf("internal error updateValue:%w", err)
	}
	return nil
}

func (rep *MetricServer) updateValues(metrics []common.Metrics) error {
	rep.logger.Debugf("updateValues len(metrics):%d", len(metrics))
	err := rep.memStorage.AddMetrics(metrics)
	if err != nil {
		return fmt.Errorf("internal error for updateValues. err:%w", err)
	}
	return nil
}
func (rep *MetricServer) UpdateJSONValue(res http.ResponseWriter, req *http.Request) {
	var metric common.Metrics
	body := req.Body
	rep.logger.Debugf("UpdateJSONValue body:%v", body)
	if err := json.NewDecoder(body).Decode(&metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	rep.logger.Debugf("UpdateJSONValue json type:%s; name: %s; value:%d; delta:%d;",
		metric.MType, metric.ID, metric.Value, metric.Delta)
	err := rep.updateValue(&metric)
	if err != nil {
		rep.logger.Debugf("Problem UpdateJSONValue json. err:%s", err)
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func (rep *MetricServer) UpdatesMetrics(res http.ResponseWriter, req *http.Request) {
	var metrics []common.Metrics
	body := req.Body
	rep.logger.Debugf("UpdatesMetrics body:%v", body)
	if err := json.NewDecoder(body).Decode(&metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := rep.updateValues(metrics)
	if err != nil {
		rep.logger.Debugf("Problem updateValues . err:%s", err)
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func (rep *MetricServer) UpdateValue(res http.ResponseWriter, req *http.Request) {
	metric, retStatus := parseURL(req, rep.logger)
	if retStatus == http.StatusOK {
		rep.logger.Debugf("update value metric:%s; delta:%s", metric, metric.Delta)
		err := rep.updateValue(&metric)
		if err != nil {
			rep.logger.Debugf("Problem UpdateJSONValue json. err:%s", err)
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
		return
	}
	res.WriteHeader(retStatus)
}

func (rep *MetricServer) getValue(metric *common.Metrics) int {
	rep.logger.Debugf("getValue type:%s; name%s; value:%d; delta:%d;", metric.MType, metric.ID, metric.Value, metric.Delta)
	retStatus := http.StatusOK
	switch metric.MType {
	case "gauge":
		val, err := rep.memStorage.GetGaugeValue(metric.ID)
		temp := float64(val)
		metric.Value = &temp
		if err != nil {
			retStatus = http.StatusNotFound
		}
	case "counter":
		val, err := rep.memStorage.GetCounterValue(metric.ID)
		temp := int64(val)
		metric.Delta = &temp
		if err != nil {
			retStatus = http.StatusNotFound
		}
	default:
		retStatus = http.StatusNotFound
		rep.logger.Warnf("unknown type:%s;", metric.MType)
	}
	return retStatus
}
func (rep *MetricServer) GetJSONValue(res http.ResponseWriter, req *http.Request) {
	var metric common.Metrics
	body := req.Body
	rep.logger.Debugf("GetJSONValue body:%v", body)
	if err := json.NewDecoder(body).Decode(&metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	rep.logger.Debugf("GetJSONValue json type:%s; name: %s; value:%d; delta:%d;",
		metric.MType, metric.ID, metric.Value, metric.Delta)
	retStatus := rep.getValue(&metric)

	resp, err := json.Marshal(metric)
	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		rep.logger.Warnf("problem with unmarshal:%w", err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(retStatus)
	_, err = res.Write(resp)
	if err != nil {
		rep.logger.Debugf("issue with write %w", err)
	}
}
func (rep *MetricServer) GetValue(res http.ResponseWriter, req *http.Request) {
	metric, retStatus := parseURL(req, rep.logger)
	if retStatus != http.StatusOK {
		res.WriteHeader(retStatus)
		return
	}

	retStatus = rep.getValue(&metric)

	res.WriteHeader(retStatus)
	_, err := res.Write([]byte(metric.GetValueMetric()))
	if err != nil {
		rep.logger.Debugf("issue for GetValue type:%s; name%s; err:%s\n", metric.MType, metric.ID, err)
	}
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
