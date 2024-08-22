package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var defaultConfig = map[string]string{
	"ServerAddress": "localhost:8080",
	"BaseURL":       "http://localhost:8080",
}

type Envs struct {
	baseURL       string `env:"BASE_URL"`
	serverAddress string `env:"SERVER_ADDRESS"`
}

type Config struct {
	ServerAddress string
	BaseURL       string
}

func (c *Config) ParseFlags() error {
	flag.StringVar(&c.ServerAddress, "a", defaultConfig["ServerAddress"], "address and port to run server")
	flag.StringVar(&c.BaseURL, "b", defaultConfig["BaseURL"], "shortLink prefix")
	flag.Parse()

	return nil
}

func (c *Config) Init() error {
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

	if c.ServerAddress == "" {
		c.ServerAddress = defaultConfig["ServerAddress"]
	}

	if c.BaseURL == "" {
		c.BaseURL = defaultConfig["BaseURL"]
	}

	return nil
}
