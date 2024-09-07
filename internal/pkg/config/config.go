package config

import (
	"github.com/caarlos0/env/v6"
)

var defaultConfig = map[string]string{
	"BaseURL":       "http://localhost:8080",
	"ServerAddress": "localhost:8080",
	"LogLevel":      "info",
}

type Envs struct {
	baseURL       string `env:"BASE_URL"`
	serverAddress string `env:"SERVER_ADDRESS"`
	logLevel      string `env:"LOG_LEVEL"`
}

type Config struct {
	BaseURL       string
	ServerAddress string
	LogLevel      string
}

func (c *Config) Init(serverAddress string, baseURL string, logLevel string) error {
	if serverAddress != "" {
		c.ServerAddress = serverAddress
	}

	if baseURL != "" {
		c.BaseURL = baseURL
	}

	if baseURL != "" {
		c.LogLevel = logLevel
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

	if envs.logLevel != "" {
		c.LogLevel = envs.logLevel
	}

	if c.ServerAddress == "" {
		c.ServerAddress = defaultConfig["ServerAddress"]
	}

	if c.BaseURL == "" {
		c.BaseURL = defaultConfig["BaseURL"]
	}

	if c.LogLevel == "" {
		c.LogLevel = defaultConfig["LogLevel"]
	}

	return nil
}
