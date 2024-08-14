package main

import (
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := "http://localhost:8080/EwHXdJfB"
		_, _ = w.Write([]byte(body))
	}

	if r.Method == http.MethodGet {
		http.Redirect(w, r, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
