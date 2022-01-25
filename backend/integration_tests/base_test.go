package integrationtest_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"
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
	ServerBaseURL string
	db            *sql.DB
}

func TestBackendTestSuite(t *testing.T) {
	backendTestSuite := &BackendTestSuite{}
	backendTestSuite.ServerBaseURL = "http://127.0.0.1:8000/api/v1"
	backendTestSuite.SetT(t)
	suite.Run(t, backendTestSuite)
}

func (suite *BackendTestSuite) SetupSuite() {
	// go func() {
	// 	grpcCMD := exec.Command("sh", "-c", "go run /__w/queue-system/queue-system/backend/integration_tests/mock_grpc/main.go")
	// 	_, err := grpcCMD.Output()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	// go func() {
	// 	backendCMD := exec.Command("sh", "-c", "go run /__w/queue-system/queue-system/backend/main.go")
	// 	_, err := backendCMD.Output()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	allDoneChan := make(chan bool)
	go func() {
		httpClient := http.Client{Timeout: 3 * time.Second}
		for {
			response, _ := httpClient.Get("http://127.0.0.1:8000" + "/api/routine/readiness")
			if response != nil && response.StatusCode == 200 {
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

	db, err := sql.Open("postgres", testConfig.POSTGRES_LOCATION())
    if err != nil {
        panic(err)
    }
	suite.db = db
}

func (suite *BackendTestSuite) TearDownSuite() {
}

func (suite *BackendTestSuite) SetupTest() {
	cmd := "/__w/queue-system/queue-system/scripts/migration_tools/migration.sh up"
	dbDown := exec.Command("sh", "-c", cmd)
	_, err := dbDown.Output()
	if err != nil {
		panic(err)
	}
}

func (suite *BackendTestSuite) TearDownTest() {
	cmd := "echo y | /__w/queue-system/queue-system/scripts/migration_tools/migration.sh down"
	dbDown := exec.Command("sh", "-c", cmd)
	_, err := dbDown.Output()
	if err != nil {
		panic(err)
	}
}
