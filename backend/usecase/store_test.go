package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB/mocks"
	"github.com/jerry0420/queue-system/backend/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpStoreTest() (
	pgDBStoreRepository *mocks.PgDBStoreRepositoryInterface,
	pgDBSignKeyRepository *mocks.PgDBSignKeyRepositoryInterface,
	storeUsecase usecase.StoreUseCaseInterface,
) {
	pgDBStoreRepository = new(mocks.PgDBStoreRepositoryInterface)
	pgDBSignKeyRepository = new(mocks.PgDBSignKeyRepositoryInterface)
	logger := logging.NewLogger([]string{}, true)
	storeUsecase = usecase.NewStoreUsecase(
		pgDBStoreRepository,
		pgDBSignKeyRepository,
		logger,
		usecase.StoreUsecaseConfig{
			Domain: "http://localhost.com",
		},
	)
	return pgDBStoreRepository, pgDBSignKeyRepository, storeUsecase
}

func TestUpdateStoreDescription(t *testing.T) {
	pgDBStoreRepository, _, storeUsecase := setUpStoreTest()

	mockStore := domain.Store{
		ID:          1,
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "description1",
		CreatedAt:   time.Now(),
		Timezone:    "Asia/Taipei",
	}
	pgDBStoreRepository.On("UpdateStore", mock.Anything, &mockStore, "description", "description1").
		Return(nil).Once()

	err := storeUsecase.UpdateStoreDescription(context.TODO(), "description1", &mockStore)
	assert.NoError(t, err)
}
