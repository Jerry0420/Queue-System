package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct{
	secretData map[string]interface{}
}

var ServerConfig Config

func init() {
	rawData, err := ioutil.ReadFile("/run/secrets/backend-secret/backend-secret")
	if err != nil {
		log.Fatalf("fail to read secret file.")
	}
	var secretData map[string]interface{}
	json.Unmarshal(rawData, &secretData)
	ServerConfig = Config{secretData: secretData}
}

func (config Config) getFromEnvVariable(key string) string {
	content := os.Getenv(key)
	if content == "" {
		// if env variable not being set properly, just exit the whole program.
		log.Fatalf("fail to validate env variable %s.", key)
	}
	return content
}

func (config Config) getFromSecretFile(key string) string {
	content, ok := config.secretData[key]
	if !ok {
		log.Fatalf("fail to get %s from secret file", key)
	}
	return content.(string)
}

type envStatus struct{ DEV, PROD, TESTING string }

var EnvStatus envStatus = envStatus{DEV: "dev", PROD: "prod", TESTING: "testing"}

// read only
func (config Config) ENV() string {
	content := config.getFromEnvVariable("ENV")
	return content
}

func (config Config) CONTEXT_TIMEOUT() time.Duration {
	content := config.getFromEnvVariable("CONTEXT_TIMEOUT")
	CONTEXT_TIMEOUT, err := strconv.Atoi(content)
	if err != nil {
		// if env variable not being set properly, just exit the whole program.
		log.Fatalf("fail to get env variable of context_timeout.")
	}
	return time.Duration(CONTEXT_TIMEOUT) * time.Second
}

func (config Config) POSTGRES_HOST() string {
	content := config.getFromSecretFile("POSTGRES_HOST")
	return content
}

func (config Config) POSTGRES_PORT() int {
	content := config.getFromSecretFile("POSTGRES_PORT")
	POSTGRES_PORT, err := strconv.Atoi(content)
	if err != nil {
		// if not set env variable properly, just exit the whole program.
		log.Fatalf("fail to get env variable of POSTGRES_PORT.")
	}
	return POSTGRES_PORT
}

func (config Config) POSTGRES_SSL() string {
	content := config.getFromEnvVariable("POSTGRES_SSL")
	return content
}

func (config Config) POSTGRES_DB() string {
	content := config.getFromSecretFile("POSTGRES_BACKEND_DB")
	return content
}

func (config Config) POSTGRES_DEV_USER() string {
	content := config.getFromEnvVariable("POSTGRES_DEV_USER")
	return content
}

func (config Config) POSTGRES_DEV_PASSWORD() string {
	content := config.getFromEnvVariable("POSTGRES_DEV_PASSWORD")
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
	content := config.getFromEnvVariable("VAULT_SERVER")
	return content
}

func (config Config) VAULT_WRAPPED_TOKEN_SERVER() string {
	content := config.getFromEnvVariable("VAULT_WRAPPED_TOKEN_SERVER")
	return content
}

func (config Config) VAULT_CRED_NAME() string {
	content := config.getFromEnvVariable("VAULT_CRED_NAME")
	return content
}

func (config Config) VAULT_ROLE_ID() string {
	content := config.getFromEnvVariable("VAULT_ROLE_ID")
	return content
}

func (config Config) VAULT_CONNECTION_CONFIG() VaultConnectionConfig {
	return VaultConnectionConfig{
		Address:             config.VAULT_SERVER(),
		WrappedTokenAddress: config.VAULT_WRAPPED_TOKEN_SERVER(),
		RoleID:              config.VAULT_ROLE_ID(),
		CredName:            config.VAULT_CRED_NAME(),
	}
}
