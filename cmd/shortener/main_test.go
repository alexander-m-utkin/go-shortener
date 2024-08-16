package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUrlHandler(t *testing.T) {
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
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			getUrlHandle(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, res.Header.Get("location"))
			}
		})
	}
}

func TestPostShortLinkHandler(t *testing.T) {
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
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			postShortLinkHandle(w, request)
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
