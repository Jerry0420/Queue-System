package integrationtest_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (suite *BackendTestSuite) Test_GetStoreInfoWithSSE() {
	httpClient := http.Client{Timeout: 3 * time.Second}

	encodedPassword := base64.StdEncoding.EncodeToString([]byte("im_password"))
	params := map[string]interface{}{
		"email":    "test@gmail.com",
		"password": encodedPassword,
		"name":     "test",
		"timezone": "Asia/Taipei",
		"queue_names": []string{
			"queue_test_1", "queue_test_2",
		},
	}
	jsonParams, _ := json.Marshal(params)
	response, _ := httpClient.Post(suite.ServerBaseURL+"/stores", "application/json", bytes.NewBuffer(jsonParams))

	suite.Equal(200, response.StatusCode)

	var decodedResponse map[string]interface{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)
	suite.Equal(1, int(decodedResponse["id"].(float64)))
	suite.Equal(2, len(decodedResponse["queues"].([]interface{})))

	params = map[string]interface{}{
		"email":    "test@gmail.com",
		"password": encodedPassword,
	}
	jsonParams, _ = json.Marshal(params)
	response, _ = httpClient.Post(suite.ServerBaseURL+"/stores/signin", "application/json", bytes.NewBuffer(jsonParams))

	suite.Equal(200, response.StatusCode)

	decodedResponse = map[string]interface{}{}
	json.NewDecoder(response.Body).Decode(&decodedResponse)
	suite.Equal(1, int(decodedResponse["id"].(float64)))

	fmt.Println(decodedResponse)
}
