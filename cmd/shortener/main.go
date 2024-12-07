package main

import (
	"flag"
	"github.com/alexander-m-utkin/go-shortener.git/internal/app"
	"github.com/alexander-m-utkin/go-shortener.git/internal/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	var serverAddressFromFlag string
	var baseURLFromFlag string
	var logLevelFromFlag string

	flag.StringVar(&serverAddressFromFlag, "a", "", "address and port to run server")
	flag.StringVar(&baseURLFromFlag, "b", "", "shortLink prefix")
	flag.StringVar(&logLevelFromFlag, "l", "", "log level")
	flag.Parse()

	err := app.Configuration.Init(serverAddressFromFlag, baseURLFromFlag, logLevelFromFlag)
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(app.Configuration.LogLevel)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(app.Configuration.ServerAddress, app.Router())
	if err != nil {
		panic(err)
	}

	logger.Log.Info("Running server", zap.String("address", app.Configuration.ServerAddress))
}
