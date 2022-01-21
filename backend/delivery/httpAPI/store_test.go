package httpAPI_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	usecaseMocks "github.com/jerry0420/queue-system/backend/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpStoreTest() (
	customerUseCase *usecaseMocks.CustomerUseCaseInterface,
	sessionUseCase *usecaseMocks.SessionUseCaseInterface,
	storeUseCase *usecaseMocks.StoreUseCaseInterface,
	integrationUseCase *usecaseMocks.IntegrationUseCaseInterface,
	httpAPIDelivery *httpAPI.HttpAPIDelivery,
	router *mux.Router,
	brokerTool *broker.Broker,
) {
	logger := logging.NewLogger([]string{}, true)
	customerUseCase = new(usecaseMocks.CustomerUseCaseInterface)
	sessionUseCase = new(usecaseMocks.SessionUseCaseInterface)
	storeUseCase = new(usecaseMocks.StoreUseCaseInterface)
	integrationUseCase = new(usecaseMocks.IntegrationUseCaseInterface)
	brokerTool = broker.NewBroker(logger)
	httpAPIDelivery = httpAPI.NewHttpAPIDelivery(
		logger,
		customerUseCase,
		sessionUseCase,
		storeUseCase,
		integrationUseCase,
		brokerTool,
		httpAPI.HttpAPIDeliveryConfig{
			StoreDuration:         time.Duration(2 * time.Second),
			TokenDuration:         time.Duration(2 * time.Second),
			PasswordTokenDuration: time.Duration(2 * time.Second),
			Domain:                "http://localhost.com",
		},
	)

	router = mux.NewRouter()
	return customerUseCase, sessionUseCase, storeUseCase, integrationUseCase, httpAPIDelivery, router, brokerTool
}

func TestOpenStore(t *testing.T) {
	_, _, storeUseCase, integrationUseCase, httpAPIDelivery, router, broker := setUpStoreTest()
	defer broker.CloseAll()

	mockStore := domain.Store{
		Name:     "name1",
		Email:    "email1",
		Password: "password1",
		Timezone: "Asia/Taipei",
	}
	mockQueues := []domain.Queue{
		{
			Name: "queue1",
		},
		{
			Name: "queue2",
		},
	}
	expectedMockStoreId := 1

	storeUseCase.On("VerifyPasswordLength", "password1").Return(nil).Once()
	storeUseCase.On("VerifyTimeZoneString", "Asia/Taipei").Return(nil).Once()
	integrationUseCase.On("CreateStore", mock.Anything, &mockStore, mockQueues).Return(nil).Run(func(args mock.Arguments) {
		store := args.Get(1).(*domain.Store)
		store.ID = expectedMockStoreId
		store.CreatedAt = time.Now()

		queues := args.Get(2).([]domain.Queue)
		queues[0].ID = 1
		queues[0].StoreID = expectedMockStoreId
		queues[1].ID = 2
		queues[1].StoreID = expectedMockStoreId
	}).Once()

	router.HandleFunc(
		httpAPI.V_1("/stores"),
		httpAPIDelivery.OpenStore,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	w := httptest.NewRecorder()

	params := map[string]interface{}{
		"name":        mockStore.Name,
		"email":       mockStore.Email,
		"password":    mockStore.Password,
		"timezone":    mockStore.Timezone,
		"queue_names": []string{mockQueues[0].Name, mockQueues[1].Name},
	}
	jsonBody, _ := json.Marshal(params)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/api/v1/stores", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var decodedResponse map[string]interface{}
	json.NewDecoder(w.Result().Body).Decode(&decodedResponse)
	assert.Equal(t, expectedMockStoreId, int(decodedResponse["id"].(float64)))
}

func TestGetStoreInfoWithSSE(t *testing.T) {
	_, _, storeUseCase, integrationUseCase, httpAPIDelivery, router, broker := setUpStoreTest()
	defer broker.CloseAll()

	expectedMockStoreDescription := "description2"
	mockStoreWithQueue := domain.StoreWithQueues{
		ID:          1,
		Name:        "name1",
		Email:       "email1",
		Description: "description1",
		CreatedAt:   time.Now(),
		Queues: []*domain.QueueWithCustomers{
			&domain.QueueWithCustomers{
				ID:   1,
				Name: "queue1",
				Customers: []*domain.Customer{
					&domain.Customer{
						ID:        1,
						Name:      "name1",
						Phone:     "phone1",
						QueueID:   1,
						Status:    domain.CustomerStatus.NORMAL,
						CreatedAt: time.Now(),
					},
					&domain.Customer{
						ID:        2,
						Name:      "name2",
						Phone:     "phone2",
						QueueID:   1,
						Status:    domain.CustomerStatus.NORMAL,
						CreatedAt: time.Now(),
					},
				},
			},
		},
	}
	mockStore := domain.Store{
		ID:          mockStoreWithQueue.ID,
		Email:       mockStoreWithQueue.Email,
		Name:        mockStoreWithQueue.Name,
		Description: expectedMockStoreDescription,
	}

	storeUseCase.On("TopicNameOfUpdateCustomer", mockStoreWithQueue.ID).Return("im_topic")
	storeUseCase.On("UpdateStoreDescription", mock.Anything, expectedMockStoreDescription, &mockStore).Return(nil).Once()
	integrationUseCase.On("GetStoreWithQueuesAndCustomersById", mock.Anything, mockStoreWithQueue.ID).Return(mockStoreWithQueue, nil)

	router.HandleFunc(
		httpAPI.V_1("/stores/{id:[0-9]+}/sse"),
		httpAPIDelivery.GetStoreInfoWithSSE,
	).Methods(http.MethodGet)

	router.HandleFunc(
		httpAPI.V_1("/stores/{id:[0-9]+}"),
		httpAPIDelivery.UpdateStoreDescription,
	).Methods(http.MethodPut)

	getContext, _ := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	putContext := context.Background()
	getW := httptest.NewRecorder()
	putW := httptest.NewRecorder()
	getDoneChan := make(chan bool)
	putDoneChan := make(chan bool)

	go func() {
		req, err := http.NewRequestWithContext(getContext, http.MethodGet, "/api/v1/stores/1/sse", nil)
		assert.NoError(t, err)
		router.ServeHTTP(getW, req)
		getDoneChan <- true
	}()

	go func() {
		time.Sleep(450 * time.Millisecond)
		params := map[string]interface{}{
			"description": expectedMockStoreDescription,
		}
		jsonBody, _ := json.Marshal(params)
		putContext = context.WithValue(putContext, domain.SignKeyTypes.NORMAL, domain.TokenClaims{
			StoreID: mockStoreWithQueue.ID,
			Email:   mockStoreWithQueue.Email,
			Name:    mockStoreWithQueue.Name,
		})
		req, err := http.NewRequestWithContext(putContext, http.MethodPut, "/api/v1/stores/1", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		router.ServeHTTP(putW, req)
		putDoneChan <- true
	}()

	<-putDoneChan
	<-getDoneChan

	re := regexp.MustCompile("data: ")
	matches := re.FindAllStringIndex(getW.Body.String(), -1)
	assert.Equal(t, 2, len(matches))
}
