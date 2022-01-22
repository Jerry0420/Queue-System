package integrationtest_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type Config struct {
	data map[string]string
}

func NewTestConfig() Config {
	secretPath := os.Getenv("SECRET_PATH")
	rawSecretData, err := ioutil.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("fail to read secret file.")
	}
	var secretData map[string]string
	json.Unmarshal(rawSecretData, &secretData)
	config := Config{data: secretData}

	envPath := os.Getenv("ENV_PATH")
	rawEnvData, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Fatalf("fail to read env file.")
	}
	var envData map[string]string
	json.Unmarshal(rawEnvData, &envData)

	for key, value := range envData {
		config.data[key] = value
	}
	return config
}

func (config Config) get(key string) string {
	content, ok := config.data[key]
	if !ok {
		log.Fatalf("fail to get %s from config", key)
	}
	return content
}

func (config Config) POSTGRES_LOCATION() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.get("POSTGRES_DEV_USER"),
		config.get("POSTGRES_DEV_PASSWORD"),
		config.get("POSTGRES_HOST"),
		config.get("POSTGRES_PORT"),
		config.get("POSTGRES_BACKEND_DB"),
		config.get("POSTGRES_SSL"),
	)
}

type BackendTestSuite struct {
	suite.Suite
}

func TestBackendTestSuite(t *testing.T) {
	backendTestSuite := &BackendTestSuite{}
	backendTestSuite.SetT(t)
	suite.Run(t, backendTestSuite)
}

func (suite *BackendTestSuite) SetupSuite() {
	go func() {
		grpcCMD := exec.Command("sh", "-c", "go run /__w/queue-system/queue-system/backend/integration_tests/mock_grpc/main.go")
		_, err := grpcCMD.Output()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		backendCMD := exec.Command("sh", "-c", "go run /__w/queue-system/queue-system/backend/main.go")
		_, err := backendCMD.Output()
		if err != nil {
			panic(err)
		}
	}()

	allDoneChan := make(chan bool)
	go func() {
		httpClient := http.Client{Timeout: 3 * time.Second}
		for {
			response, _ := httpClient.Get("http://127.0.0.1:8000" + "/api/routine/readiness")
			if response != nil && response.StatusCode == 200 {
				bodyBytes, _ := io.ReadAll(response.Body)
				fmt.Println(string(bodyBytes))
				break
			}
		}
		allDoneChan <- true
	}()
	<-allDoneChan

	testConfig := NewTestConfig()
	os.Setenv("POSTGRES_MIGRATION_USER", testConfig.get("POSTGRES_DEV_USER"))
	os.Setenv("POSTGRES_MIGRATION_PASSWORD", testConfig.get("POSTGRES_DEV_PASSWORD"))
	os.Setenv("POSTGRES_HOST", testConfig.get("POSTGRES_HOST"))
	os.Setenv("POSTGRES_PORT", testConfig.get("POSTGRES_PORT"))
	os.Setenv("POSTGRES_BACKEND_DB", testConfig.get("POSTGRES_BACKEND_DB"))
}

func (suite *BackendTestSuite) TearDownSuite() {
}

func (suite *BackendTestSuite) SetupTest() {
}

func (suite *BackendTestSuite) TearDownTest() {
}

func (suite *BackendTestSuite) Test_aaa() {
	suite.Equal(1, 1)
}

func (suite *BackendTestSuite) Test_bbb() {
	suite.Equal(1, 1)
}