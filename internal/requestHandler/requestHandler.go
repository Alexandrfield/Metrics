package requesthandler

import (
	"html/template"
	"net/http"
	"strings"

	"log"
)

type MetricsStorage interface {
	SetValue(metricType string, metricName string, metricValue string) error
	GetValue(metricType string, metricName string) (string, error)
	GetAllValue() ([]string, error)
}

type MetricServer struct {
	memStorage MetricsStorage
}

func CreateHandlerRepository(stor MetricsStorage) *MetricServer {
	return &MetricServer{memStorage: stor}
}

func parseURL(req *http.Request) ([]string, int) {
	url := strings.Split(req.URL.String(), "/")
	// expected format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain
	if url[1] == "update" && len(url) != 5 {
		return []string{}, http.StatusNotFound
	}
	//expected format http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
	if url[1] == "value" && len(url) < 4 {
		return []string{}, http.StatusBadRequest
	}
	return url, http.StatusOK
}
func (rep *MetricServer) DefaultAnswer(res http.ResponseWriter, req *http.Request) {
	log.Printf("defaultAnswer. req:%v;res.WriteHeader::%d\n", req, http.StatusNotImplemented)
	res.WriteHeader(http.StatusNotImplemented)
}

func (rep *MetricServer) UpdateValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusMethodNotAllowed

	var url []string
	url, statusH = parseURL(req)
	if statusH == http.StatusOK {
		err := rep.memStorage.SetValue(url[3], url[3], url[4])
		if err != nil {
			log.Printf("issue for updateValue type:%s; name%s; value:%s; err:%s\n", url[3], url[3], url[4], err)
			statusH = http.StatusBadRequest
		}
	}

	log.Printf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}
func (rep *MetricServer) GetValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusMethodNotAllowed

	var url []string
	url, statusH = parseURL(req)
	if statusH == http.StatusOK {
		val, err := rep.memStorage.GetValue(url[2], url[3])
		if err != nil {
			res.WriteHeader(statusH)
			res.Write([]byte(val))
			return
		} else {
			log.Printf("issue for GetValue type:%s; name%s; err:%s\n", url[3], url[3], err)
			statusH = http.StatusNotFound
		}
	}

	log.Printf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}

func (rep *MetricServer) GetAllData(res http.ResponseWriter, req *http.Request) {
	allValues, err := rep.memStorage.GetAllValue()
	if err == nil {
		log.Printf("issue for GetAllData. err:%s\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	all_values_template := `
<html> 
   <head> 
   </head> 
   <body> 
	all metrics: 
		{{ range .}}{{.}}, {{ end }}
   </body> 
</html>
`
	ready_template, err := template.New("templ").Parse(all_values_template)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	ready_template.Execute(res, allValues)
}
