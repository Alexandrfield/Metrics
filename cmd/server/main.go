package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	handl "github.com/Alexandrfield/Metrics/internal/requestHandler"
)

func parseURL(req *http.Request) ([]string, int) {
	url := strings.Split(req.URL.String(), "/")
	fmt.Printf("parse len():%d, url:%v\n", len(url), url)
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
func defaultAnswer(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("defaultAnswer. req:%v;\n", req)
	fmt.Printf("defaultAnswer: res.WriteHeader:%d\n", http.StatusNotImplemented)
	res.WriteHeader(http.StatusNotImplemented)
}

func updateValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusMethodNotAllowed
	fmt.Printf("req:%v; req.Method:%s\n", req, req.Method)

	var url []string
	url, statusH = parseURL(req)
	if statusH == http.StatusOK {
		fmt.Printf("try save url\n")
		if !handl.HandleRequest(url) {
			statusH = http.StatusBadRequest
		}
	}

	fmt.Printf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}
func getValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusMethodNotAllowed

	var url []string
	url, statusH = parseURL(req)
	if statusH == http.StatusOK {
		val, st := handl.HandleGetValue(url)
		if st {
			res.WriteHeader(statusH)
			res.Write([]byte(val))
			return
		} else {
			statusH = http.StatusNotFound
		}
	}

	fmt.Printf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}

func getAllData(res http.ResponseWriter, req *http.Request) {
	allValues := handl.HandleAllValue()
	page := `
<html> 
   <head> 
   </head> 
   <body> 
`
	for i := 0; i < len(allValues); i++ {
		page += fmt.Sprintf(`<h3>%s   </h3>`, allValues[i])
	}
	page += `
   </body> 
</html>
`
	res.Header().Set("content-type", "Content-Type: text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(page))
}
func main() {
	router := chi.NewRouter()

	router.Post(`/update/gauge/`, updateValue)
	router.Post(`/update/counter/`, updateValue)
	router.Post(`/update/`, defaultAnswer)
	router.Get(`/value/`, getValue)
	router.Get(`/`, getAllData)

	err := http.ListenAndServe(`127.0.0.1:8080`, router)
	if err != nil {
		panic(err)
	}
}
