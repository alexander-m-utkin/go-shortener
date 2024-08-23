package main

import (
	"flag"
	"github.com/alexander-m-utkin/go-shortener.git/internal/app"
	"log"
	"net/http"
)

func main() {
	var serverAddressFromFlag string
	var baseURLFromFlag string

	flag.StringVar(&serverAddressFromFlag, "a", "", "address and port to run server")
	flag.StringVar(&baseURLFromFlag, "b", "", "shortLink prefix")
	flag.Parse()

	err := app.Configuration.Init(serverAddressFromFlag, baseURLFromFlag)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(app.Configuration.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}
}
