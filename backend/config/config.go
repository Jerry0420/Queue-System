package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
)

type Config struct {
	logger logging.LoggerTool
}

func NewConfig(logger logging.LoggerTool) Config {
	config := Config{logger}
	return config
}

func (config Config) validate(content string) string {
	if content == "" {
		// if env variable not being set properly, just exit the whole program.
		config.logger.FATALf("fail to validate env variable.")
	}
	return content
}

type envStatus struct {DEV, PROD, TESTING string}
var EnvStatus envStatus = envStatus{DEV: "dev", PROD: "prod", TESTING: "testing"}

// read only
func (config Config) ENV() string {
	content := config.validate(os.Getenv("ENV"))
	return content
}

func (config Config) CONTEXT_TIMEOUT() time.Duration {
	content := config.validate(os.Getenv("CONTEXT_TIMEOUT"))
	CONTEXT_TIMEOUT, err := strconv.Atoi(content)
	if err != nil {
		// if env variable not being set properly, just exit the whole program.
		config.logger.FATALf("fail to get env variable of context_timeout.")
	}
	return time.Duration(CONTEXT_TIMEOUT) * time.Second
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
		config.logger.FATALf("fail to get env variable of POSTGRES_PORT.")
	}
	return POSTGRES_PORT
}

func (config Config) POSTGRES_SSL() string {
	content := config.validate(os.Getenv("POSTGRES_SSL"))
	return content
}

func (config Config) POSTGRES_DB() string {
	content := config.validate(os.Getenv("POSTGRES_BACKEND_DB"))
	return content
}

func (config Config) POSTGRES_DEV_USER() string {
	content := config.validate(os.Getenv("POSTGRES_DEV_USER"))
	return content
}

func (config Config) POSTGRES_DEV_PASSWORD() string {
	content := config.validate(os.Getenv("POSTGRES_DEV_PASSWORD"))
	return content
}

func (config Config) POSTGRES_LOCATION() string {
	return fmt.Sprintf("%s:%d/%s?sslmode=%s", 
		config.POSTGRES_HOST(), 
		config.POSTGRES_PORT(), 
		config.POSTGRES_DB(),
		config.POSTGRES_SSL(),
    )
}

func (config Config) VAULT_SERVER() string {
	content := config.validate(os.Getenv("VAULT_SERVER"))
	return content
}

func (config Config) VAULT_WRAPPED_TOKEN_SERVER() string {
	content := config.validate(os.Getenv("VAULT_WRAPPED_TOKEN_SERVER"))
	return content
}

func (config Config) VAULT_CRED_NAME() string {
	content := config.validate(os.Getenv("VAULT_CRED_NAME"))
	return content
}

func (config Config) VAULT_ROLE_ID() string {
	content := config.validate(os.Getenv("VAULT_ROLE_ID"))
	return content
}

func (config Config) VAULT_CONNECTION_CONFIG() VaultConnectionConfig {
	return VaultConnectionConfig{
		Address: config.VAULT_SERVER(),
		WrappedTokenAddress: config.VAULT_WRAPPED_TOKEN_SERVER(),
		RoleID: config.VAULT_ROLE_ID(),
		CredName: config.VAULT_CRED_NAME(),
	}
}