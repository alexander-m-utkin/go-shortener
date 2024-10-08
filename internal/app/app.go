package app

import (
	"github.com/alexander-m-utkin/go-shortener.git/internal/pkg/config"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
)

var Configuration config.Config

var GlobalStorage = map[string]string{
	"EwHXdJfB": "https://practicum.yandex.ru/",
	"NQmfnwrt": "https://practicum.yandex.ru/123/",
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
	// если url уже есть в хранилище то просто вернем его id,
	// если url нет в хранилище - сгенерируем новый id и сохраним под ним url.
	// Тут перебор map KeyForValue это временный вариант, конечно так не надо делать.
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

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", GetURLHandle)
	r.Post("/", PostShortLinkHandle)
	return r
}
