package config

import "go-base-structure/pkg/env"

// Config is our configuration about the server
type Config struct {
	HTTPPort string
	Domain   string
	Debug    bool
}

func NewConfig() *Config {
	HTTPPort := env.GetEnvOrDefaultString("HTTP_PORT", "8000")
	Domain := env.GetEnvOrDefaultString("DOMAIN", "localhost")
	Debug := env.GetEnvOrDefaultBool("DEBUG", true)
	return &Config{
		HTTPPort: HTTPPort,
		Domain:   Domain,
		Debug:    Debug,
	}
}
