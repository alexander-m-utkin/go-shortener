package main

import (
	"github.com/alexander-m-utkin/go-shortener.git/internal/app"
	"log"
	"net/http"
)

func main() {
	err := app.Configuration.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(app.Configuration.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
