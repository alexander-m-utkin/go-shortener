package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Envs struct {
	baseURL       string `env:"BASE_URL"`
	serverAddress string `env:"SERVER_ADDRESS"`
}

type Config struct {
	ServerAddress string
	BaseURL       string
	isInitiated   bool
}

func (c *Config) Init() error {
	if !c.isInitiated {
		flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
		flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "shortLink prefix")
		flag.Parse()

		var envs Envs
		err := env.Parse(&envs)
		if err != nil {
			return err
		}

		if envs.serverAddress != "" {
			c.ServerAddress = envs.serverAddress
		}

		if envs.baseURL != "" {
			c.BaseURL = envs.baseURL
		}

		c.isInitiated = true

	}
	return nil
}
