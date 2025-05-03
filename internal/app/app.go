package app

import (
	"encoding/json"
	"github.com/alexander-m-utkin/go-shortener.git/internal/pkg/compress"
	"github.com/alexander-m-utkin/go-shortener.git/internal/pkg/config"
	"github.com/alexander-m-utkin/go-shortener.git/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
)

var Configuration config.Config

var GlobalStorage = map[string]string{
	"EwHXdJfB": "https://practicum.yandex.ru/",
	"NQmfnwrt": "https://practicum.yandex.ru/123/",
	"OWVHDnHE": "http://h02ifaoltlwna.yandex/c0sdvtdv/jlhiim6aqxo6pf",
	"iWMRzXOr": "http://viwwcdahrwag2h.ru/i6nqys7ksj/fbvx6fpddpowr",
	"KzZvYuMr": "http://n3tbmjo.biz/cp4pcntmrwpu/lofqr",
}

func RandString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func KeyForValue(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}

func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if url, ok := GlobalStorage[id]; ok {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func PostShortLinkHandle(w http.ResponseWriter, r *http.Request) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request rBody", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rBodyString := string(rBody)

	var id string
	if foundKey, isFound := KeyForValue(GlobalStorage, rBodyString); isFound {
		id = foundKey
	} else {
		id = RandString(8)
		GlobalStorage[id] = string(rBody)
	}

	shortLink := Configuration.BaseURL + "/" + id

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortLink))
}

type PostShortenRequest struct {
	URL string `json:"url"`
}

type PostShortenResponse struct {
	Result string `json:"result"`
}

func PostShortenHandle(w http.ResponseWriter, r *http.Request) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request rBody", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var body PostShortenRequest

	if err = json.Unmarshal(rBody, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id string
	if foundKey, isFound := KeyForValue(GlobalStorage, body.URL); isFound {
		id = foundKey
	} else {
		id = RandString(8)
		GlobalStorage[id] = body.URL
	}

	shortLink := Configuration.BaseURL + "/" + id

	wBody := PostShortenResponse{
		Result: shortLink,
	}

	wBodyJSON, err := json.Marshal(wBody)
	if err != nil {
		http.Error(w, "Failed marshal string", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(wBodyJSON)
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Use(compress.ZlibHandle)
	r.Use(logger.RequestLogger)

	r.Post("/api/shorten", PostShortenHandle)
	r.Get("/{id}", GetURLHandle)
	r.Post("/", PostShortLinkHandle)
	return r
}
