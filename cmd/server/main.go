package main

import (
	"fmt"
	"net/http"
	"strings"

	handl "github.com/Alexandrfield/Metrics/internal/requestHandler"
)

func parseURL(req *http.Request) ([]string, int) {
	url := strings.Split(req.URL.String(), "/")
	// expected format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain
	if len(url) != 5 {
		return []string{}, http.StatusNotFound
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
	if req.Method == http.MethodPost {
		var url []string
		url, statusH = parseURL(req)
		fmt.Printf("parse st:%d, url:%v\n", statusH, url)
		if statusH == http.StatusOK {
			fmt.Printf("try save url\n")
			if !handl.HandleRequest(url) {
				statusH = http.StatusBadRequest
			}
		}
	}
	fmt.Printf("res.WriteHeader:%d\n", statusH)
	res.WriteHeader(statusH)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/gauge/`, updateValue)
	mux.HandleFunc(`/update/counter/`, updateValue)
	mux.HandleFunc(`/update/`, defaultAnswer)

	err := http.ListenAndServe(`127.0.0.1:8080`, mux)
	if err != nil {
		panic(err)
	}
}
