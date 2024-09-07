package main

import (
	"bytes"
	"github.com/alexander-m-utkin/go-shortener.git/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(ts *httptest.Server, method,
	path string, body string) *resty.Response {

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	req := client.R()

	if body != "" {
		client.R().SetBody(body)
	}

	resp, err := req.Execute(method, ts.URL+path)

	if err != nil {
		// Проверяем, если ошибка возникла НЕ из-за отсутствия следования за редиректами
		if resp != nil && resp.StatusCode() != http.StatusTemporaryRedirect {
			panic(err)
		}
	}

	return resp
}

func TestRouter(t *testing.T) {
	err := app.Configuration.Init("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewServer(app.Router())
	defer ts.Close()

	type want struct {
		body     string
		code     int
		location string
	}

	var testTable = []struct {
		method string
		url    string
		want   want
		body   string
	}{
		{method: http.MethodGet, url: "/EwHXdJfB", want: want{code: http.StatusTemporaryRedirect, location: "https://practicum.yandex.ru/"}},
		{method: http.MethodGet, url: "/33333333", want: want{code: http.StatusBadRequest}},
		{method: http.MethodPost, url: "/", body: "https://practicum.yandex.ru/", want: want{code: http.StatusCreated}},
		{method: http.MethodPut, url: "/", body: "https://practicum.yandex.ru/", want: want{code: http.StatusMethodNotAllowed}},
	}
	for _, tt := range testTable {
		resp := testRequest(ts, tt.method, tt.url, tt.body)

		assert.Equal(t, tt.want.code, resp.StatusCode())

		if tt.want.location != "" {
			assert.Equal(t, tt.want.location, resp.Header().Get("location"))
		}

		if tt.want.body != "" {
			assert.Equal(t, tt.body, string(resp.Body()))
		}
	}
}

func TestGetUrlHandler(t *testing.T) {
	err := app.Configuration.Init("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	type want struct {
		body        string
		code        int
		contentType string
		location    string
	}

	tests := []struct {
		body   string
		method string
		name   string
		path   string
		want   want
	}{
		{
			method: http.MethodGet,
			name:   "GET /EwHXdJfB",
			path:   "/EwHXdJfB",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				body:        "",
				location:    "https://practicum.yandex.ru/",
			},
		},
		{
			method: http.MethodPost,
			name:   "POST /EwHXdJfB",
			path:   "/EwHXdJfB",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			r := chi.NewRouter()
			r.Get("/{id}", app.GetURLHandle)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, res.Header.Get("location"))
			}
		})
	}
}

func TestPostShortLinkHandler(t *testing.T) {
	err := app.Configuration.Init("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	type want struct {
		body string
		code int
	}

	tests := []struct {
		body   string
		method string
		name   string
		path   string
		want   want
	}{
		{
			name:   "POST /",
			body:   "https://practicum.yandex.ru/",
			method: http.MethodPost,
			path:   "/",
			want: want{
				body: "http://localhost:8080/EwHXdJfB",
				code: http.StatusCreated,
			},
		},
		{
			name:   "PUT /",
			body:   "https://practicum.yandex.ru/",
			method: http.MethodPut,
			path:   "/",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			request.Host = "localhost:8080"

			r := chi.NewRouter()
			r.Post("/", app.PostShortLinkHandle)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.body != "" {
				bodyBytes, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("Failed to read body: %v", err)
				}
				defer res.Body.Close()

				bodyString := string(bodyBytes)

				assert.Equal(t, tt.want.body, bodyString)
			}
		})
	}
}
