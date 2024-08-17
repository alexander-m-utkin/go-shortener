package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Envs struct {
	baseUrl       string `env:"BASE_URL"`
	serverAddress string `env:"SERVER_ADDRESS"`
}

type Config struct {
	ServerAddress string
	BaseUrl       string
	isInitiated   bool
}

func (c *Config) Init() {
	if c.isInitiated != true {
		flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
		flag.StringVar(&c.BaseUrl, "b", "http://localhost:8080", "shortLink prefix")
		flag.Parse()

		var envs Envs
		err := env.Parse(&envs)
		if err != nil {
			log.Fatal(err)
		}

		if envs.serverAddress != "" {
			c.ServerAddress = envs.serverAddress
		}

		if envs.baseUrl != "" {
			c.BaseUrl = envs.baseUrl
		}

		c.isInitiated = true
	}
}
