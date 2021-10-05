package config

import (
	"os"
	"strconv"
)

type Config struct {}

func NewConfig() Config {
	config := Config{}
	return config
}

func (config Config) validate(content string) string {
	if content == "" {
		// if env variable not being set properly, just exit the whole program.
		os.Exit(1)
	}
	return content
}

// read only
func (config Config) ENV() string {
	content := config.validate(os.Getenv("ENV"))
	return content
}

func (config Config) CONTEXT_TIMEOUT() int {
	content := config.validate(os.Getenv("CONTEXT_TIMEOUT"))
	CONTEXT_TIMEOUT, err := strconv.Atoi(content)
	if err != nil {
		// if env variable not being set properly, just exit the whole program.
		os.Exit(1)
	}
	return CONTEXT_TIMEOUT
}

func (config Config) POSTGRES_HOST() string {
	content := config.validate(os.Getenv("POSTGRES_HOST"))
	return content
}

func (config Config) POSTGRES_PORT() int {
	content := config.validate(os.Getenv("POSTGRES_PORT"))
	POSTGRES_PORT, err := strconv.Atoi(content)
	if err != nil {
		// if not set env variable properly, just exit the whole program.
		os.Exit(1)
	}
	return POSTGRES_PORT
}

func (config Config) POSTGRES_DB() string {
	content := config.validate(os.Getenv("POSTGRES_DB"))
	return content
}

func (config Config) POSTGRES_USER() string {
	content := config.validate(os.Getenv("POSTGRES_USER"))
	return content
}

func (config Config) POSTGRES_PASSWORD() string {
	content := config.validate(os.Getenv("POSTGRES_PASSWORD"))
	return content
}