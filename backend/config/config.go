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

type Config struct {
	data map[string]string
}

var ServerConfig Config

func init() {
	secretPath := "/run/secrets/backend-secret/.backend-secret"
	if os.Getenv("SECRET_PATH") != "" {
		secretPath = os.Getenv("SECRET_PATH")
	}
	rawSecretData, err := ioutil.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("fail to read secret file.")
	}
	var secretData map[string]string
	json.Unmarshal(rawSecretData, &secretData)
	ServerConfig = Config{data: secretData}

	envPath := "/etc/config/.env"
	if os.Getenv("ENV_PATH") != "" {
		envPath = os.Getenv("ENV_PATH")
	}
	rawEnvData, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Fatalf("fail to read env file.")
	}
	var envData map[string]string
	json.Unmarshal(rawEnvData, &envData)

	for key, value := range envData {
		ServerConfig.data[key] = value
	}
}

func (config Config) get(key string) string {
	content, ok := config.data[key]
	if !ok {
		log.Fatalf("fail to get %s from config", key)
	}
	return content
}

type envState struct{ DEV, PROD, TESTING string }

var EnvState envState = envState{DEV: "dev", PROD: "prod", TESTING: "testing"}

// read only
func (config Config) ENV() string {
	content := config.get("ENV")
	return content
}

func (config Config) DOMAIN() string {
	content := config.get("DOMAIN")
	return content
}

func (config Config) CONTEXT_TIMEOUT() time.Duration {
	content := config.get("CONTEXT_TIMEOUT")
	CONTEXT_TIMEOUT, err := strconv.Atoi(content)
	if err != nil {
		// if env variable not being set properly, just exit the whole program.
		log.Fatalf("fail to get env variable of context_timeout.")
	}
	return time.Duration(CONTEXT_TIMEOUT) * time.Second
}

func (config Config) STOREDURATION() time.Duration {
	return 24 * time.Hour
}

func (config Config) TOKENDURATION() time.Duration {
	return 6 * time.Hour
}

func (config Config) PASSWORDTOKENDURATION() time.Duration {
	return 5 * time.Minute
}

func (config Config) GRPCREPLICACOUNT() int {
	return 2
}

func (config Config) TOKEN_SIGN_KEY() string {
	content := config.get("TOKEN_SIGN_KEY")
	return content
}

func (config Config) POSTGRES_HOST() string {
	content := config.get("POSTGRES_HOST")
	return content
}

func (config Config) POSTGRES_PORT() int {
	content := config.get("POSTGRES_PORT")
	POSTGRES_PORT, err := strconv.Atoi(content)
	if err != nil {
		// if not set env variable properly, just exit the whole program.
		log.Fatalf("fail to get env variable of POSTGRES_PORT.")
	}
	return POSTGRES_PORT
}

func (config Config) POSTGRES_SSL() string {
	content := config.get("POSTGRES_SSL")
	return content
}

func (config Config) POSTGRES_DB() string {
	content := config.get("POSTGRES_BACKEND_DB")
	return content
}

func (config Config) POSTGRES_DEV_USER() string {
	content := config.get("POSTGRES_DEV_USER")
	return content
}

func (config Config) POSTGRES_DEV_PASSWORD() string {
	content := config.get("POSTGRES_DEV_PASSWORD")
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
	content := config.get("VAULT_SERVER")
	return content
}

func (config Config) VAULT_WRAPPED_TOKEN_SERVER() string {
	content := config.get("VAULT_WRAPPED_TOKEN_SERVER")
	return content
}

func (config Config) VAULT_CRED_NAME() string {
	content := config.get("VAULT_CRED_NAME")
	return content
}

func (config Config) VAULT_ROLE_ID() string {
	content := config.get("VAULT_ROLE_ID")
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

func (config Config) GRPC_HOST() string {
	content := config.get("GRPC_HOST")
	return content
}

func (config Config) CA_CRT() string {
	content := config.get("CA_CRT")
	return content
}
