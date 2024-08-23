package config

import (
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

func (c *Config) Init(serverAddress string, baseURL string) error {
	if serverAddress != "" {
		c.ServerAddress = serverAddress
	}

	if baseURL != "" {
		c.BaseURL = baseURL
	}

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
