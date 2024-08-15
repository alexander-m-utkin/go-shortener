package main

import (
	"net/http"
)

func getUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect)
}

func postShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	body := "http://localhost:8080/EwHXdJfB"
	_, _ = w.Write([]byte(body))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, getUrlHandler)
	mux.HandleFunc(`/`, postShortLinkHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
