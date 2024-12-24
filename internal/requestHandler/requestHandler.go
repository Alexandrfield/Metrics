package requesthandler

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Alexandrfield/Metrics/internal/common"
)

type MetricsStorage interface {
	SetValue(metricType string, metricName string, metricValue string) error
	GetValue(metricType string, metricName string) (string, error)
	GetAllValue() ([]string, error)
}

type MetricServer struct {
	logger     common.Loger
	memStorage MetricsStorage
}

func CreateHandlerRepository(stor MetricsStorage) *MetricServer {
	return &MetricServer{memStorage: stor}
}

func parseURL(req *http.Request) ([]string, int) {
	url := strings.Split(req.URL.String(), "/")
	// expected format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>,
	// Content-Type: text/plain
	if url[1] == "update" && len(url) != 5 {
		return []string{}, http.StatusNotFound
	}
	// expected format http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	if url[1] == "value" && len(url) < 4 {
		return []string{}, http.StatusBadRequest
	}
	return url, http.StatusOK
}
func (rep *MetricServer) DefaultAnswer(res http.ResponseWriter, req *http.Request) {
	rep.logger.Debugf("defaultAnswer. req:%v;res.WriteHeader::%d\n", req, http.StatusNotImplemented)
	res.WriteHeader(http.StatusNotImplemented)
}

func (rep *MetricServer) UpdateValue(res http.ResponseWriter, req *http.Request) {
	url, statusH := parseURL(req)
	if statusH == http.StatusOK {
		err := rep.memStorage.SetValue(url[2], url[3], url[4])
		rep.logger.Debugf("setValue type:%s; name%s; value:%s; err:%s\n", url[2], url[3], url[4], err)
		if err != nil {
			rep.logger.Debugf("issue for updateValue type:%s; name%s; value:%s; err:%s\n", url[2], url[3], url[4], err)
			statusH = http.StatusBadRequest
		}
	}

	rep.logger.Debugf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}
func (rep *MetricServer) GetValue(res http.ResponseWriter, req *http.Request) {
	url, statusH := parseURL(req)
	if statusH == http.StatusOK {
		rep.logger.Debugf("GetValue(url[2], url[3])> %s, %s\n", url[2], url[3])
		val, err := rep.memStorage.GetValue(url[2], url[3])
		if err != nil {
			rep.logger.Debugf("issue for res.Write([]byte(val)); err:%s\n", err)
			statusH = http.StatusNotFound
		} else {
			res.WriteHeader(statusH)
			_, err = res.Write([]byte(val))
			if err != nil {
				rep.logger.Debugf("issue for GetValue type:%s; name%s; err:%s\n", url[2], url[3], err)
			}
			return
		}
	}

	rep.logger.Debugf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
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
