package main

import (
	"net/http"
	"strings"

	handl "github.com/Alexandrfield/Metrics/internal/requestHandler"
)

func parseUrl(req *http.Request) ([]string, int) {
	url := strings.Split(req.URL.String(), "/")
	// expected format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain
	if len(url) != 5 {
		return []string{}, http.StatusNotFound
	}
	return url, http.StatusOK
}
func defaultAnswer(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}

func updateValue(res http.ResponseWriter, req *http.Request) {
	statusH := http.StatusMethodNotAllowed
	if req.Method != http.MethodPost {
		url, st := parseUrl(req)
		res.WriteHeader(st)
		if st == http.StatusOK {
			handl.HandleRequest(url)
		}
	}
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