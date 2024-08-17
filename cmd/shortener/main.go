package main

import (
	"github.com/alexander-m-utkin/go-shortener.git/config"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
)

var configuration config.Config

var globalStorage = map[string]string{
	"EwHXdJfB": "https://practicum.yandex.ru/",
	"NQmfnwrt": "https://practicum.yandex.ru/123/",
}

func randString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func keyForValue(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}

func getUrlHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if url, ok := globalStorage[id]; ok {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func postShortLinkHandle(w http.ResponseWriter, r *http.Request) {
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
	// Тут перебор map keyForValue это временный вариант, конечно так не надо делать.
	if foundKey, isFound := keyForValue(globalStorage, rBodyString); isFound {
		id = foundKey
	} else {
		id = randString(8)
		globalStorage[id] = string(rBody)
	}

	shortLink := configuration.BaseUrl + "/" + id

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortLink))
}

func AppRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", getUrlHandle)
	r.Post("/", postShortLinkHandle)
	return r
}

func main() {
	configuration.Init()

	err := http.ListenAndServe(configuration.ServerAddress, AppRouter())
	if err != nil {
		panic(err)
	}
}
