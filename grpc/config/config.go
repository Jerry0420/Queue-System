package config

import (
	"encoding/json"
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
	path := "/run/secrets/grpc-secret/.grpc-secret"
	if os.Getenv("GRPC_SECRET_PATH") != "" {
		path = os.Getenv("GRPC_SECRET_PATH")
	}
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("fail to read secret file.")
	}
	var secretData map[string]string
	json.Unmarshal(rawData, &secretData)
	ServerConfig = Config{data: secretData}

	envPath := "/etc/config/.env"
	if os.Getenv("GRPC_ENV_PATH") != "" {
		envPath = os.Getenv("GRPC_ENV_PATH")
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

type envStatus struct{ DEV, PROD, TESTING string }

var EnvStatus envStatus = envStatus{DEV: "dev", PROD: "prod", TESTING: "testing"}

// read only
func (config Config) ENV() string {
	content := config.get("ENV")
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

func (config Config) EMAIL_FROM() string {
	content := config.get("EMAIL_FROM")
	return content
}

func (config Config) EMAIL_SERVER() string {
	content := config.get("EMAIL_SERVER")
	return content
}

func (config Config) EMAIL_USERNAME() string {
	content := config.get("EMAIL_USERNAME")
	return content
}

func (config Config) EMAIL_PASSWORD() string {
	content := config.get("EMAIL_PASSWORD")
	return content
}

func (config Config) SERVER_CRT() string {
	content := config.get("SERVER_CRT")
	return content
}

func (config Config) SERVER_KEY() string {
	content := config.get("SERVER_KEY")
	return content
}
