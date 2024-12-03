package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateValue(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "/update/counter/test/1",
			want: want{
				code:        200,
				contentType: "plain/text",
			},
		},
		{
			name:    "positive test #2",
			request: "/update/gauge/test/4.4",
			want: want{
				code:        200,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #2",
			request: "/update/counter/5",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
		{
			name:    "negativ test #3",
			request: "/update/gauge/5",
			want: want{
				code:        http.StatusNotFound,
				contentType: "plain/text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			updateValue(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}
func TestDefaultAnswer(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/update/unknown/test/1", nil)
	w := httptest.NewRecorder()

	defaultAnswer(w, request)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusNotImplemented, result.StatusCode)

}

func TestParserURL(t *testing.T) {
	tests := []struct {
		name      string
		request   string
		resStatus int
	}{
		{
			name:      "positive test #1",
			request:   "/update/counter/test/1",
			resStatus: http.StatusOK,
		},
		{
			name:      "negative test #1",
			request:   "/update/counter/1",
			resStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			updateValue(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.resStatus, result.StatusCode)
		})
	}

}
