package main

import (
	"github.com/alexander-m-utkin/go-shortener.git/internal/app"
	"net/http"
)

func main() {
	app.Configuration.Init()

	err := http.ListenAndServe(app.Configuration.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
