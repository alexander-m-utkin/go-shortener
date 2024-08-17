package config

import "flag"

type Config struct {
	A           string
	B           string
	isInitiated bool
}

func (c *Config) Init() {
	if c.isInitiated != true {
		flag.StringVar(&c.A, "a", "localhost:8080", "address and port to run server")
		flag.StringVar(&c.B, "b", "http://localhost:8080", "shortLink prefix")

		flag.Parse()

		c.isInitiated = true
	}
}
