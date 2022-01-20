package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	grpcServicesMocks "github.com/jerry0420/queue-system/backend/repository/grpcServices/mocks"
	pgDBMocks "github.com/jerry0420/queue-system/backend/repository/pgDB/mocks"
	"github.com/jerry0420/queue-system/backend/usecase"
	usecaseMocks "github.com/jerry0420/queue-system/backend/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpIntegrationTest() (
	pgDBTx *pgDBMocks.PgDBTxInterface,
	pgDBStoreRepository *pgDBMocks.PgDBStoreRepositoryInterface,
	pgDBSessionRepository *pgDBMocks.PgDBSessionRepositoryInterface,
	pgDBCustomerRepository *pgDBMocks.PgDBCustomerRepositoryInterface,
	pgDBQueueRepository *pgDBMocks.PgDBQueueRepositoryInterface,
	grpcServicesRepository *grpcServicesMocks.GrpcServicesRepositoryInterface,
	storeUseCase *usecaseMocks.StoreUseCaseInterface,
	integrationUsecase usecase.IntegrationUseCaseInterface,
	pgDB *pgDBMocks.PgDBInterface,
) {
	pgDBTx = new(pgDBMocks.PgDBTxInterface)
	pgDBStoreRepository = new(pgDBMocks.PgDBStoreRepositoryInterface)
	pgDBSessionRepository = new(pgDBMocks.PgDBSessionRepositoryInterface)
	pgDBCustomerRepository = new(pgDBMocks.PgDBCustomerRepositoryInterface)
	pgDBQueueRepository = new(pgDBMocks.PgDBQueueRepositoryInterface)
	grpcServicesRepository = new(grpcServicesMocks.GrpcServicesRepositoryInterface)
	storeUseCase = new(usecaseMocks.StoreUseCaseInterface)
	logger := logging.NewLogger([]string{}, true)

	integrationUsecase = usecase.NewIntegrationUsecase(
		pgDBTx,
		pgDBStoreRepository,
		pgDBSessionRepository,
		pgDBCustomerRepository,
		pgDBQueueRepository,
		grpcServicesRepository,
		storeUseCase,
		logger,
		usecase.IntegrationUsecaseConfig{
			StoreDuration:         time.Duration(5 * time.Second),
			TokenDuration:         time.Duration(5 * time.Second),
			PasswordTokenDuration: time.Duration(5 * time.Second),
			GrpcReplicaCount:      3,
		},
	)
	pgDB = new(pgDBMocks.PgDBInterface)
	return pgDBTx, pgDBStoreRepository, pgDBSessionRepository, pgDBCustomerRepository, pgDBQueueRepository, grpcServicesRepository, storeUseCase, integrationUsecase, pgDB
}

func TestCreateStore(t *testing.T) {
	pgDBTx, pgDBStoreRepository, _, _, pgDBQueueRepository, _, storeUseCase, integrationUsecase, pgDB := setUpIntegrationTest()
	mockStore := domain.Store{
		ID:          1,
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "description1",
		Timezone: "Asia/Taipei",
	}
	mockQueues := []domain.Queue{
		{
			ID:      1,
			Name:    "queue1",
			StoreID: 1,
		},
		{
			ID:      2,
			Name:    "queue2",
			StoreID: 1,
		},
	}

	storeUseCase.On("EncryptPassword", "password1").Return("encryptPassword1", nil).Once()
	pgDBTx.On("BeginTx").Return(pgDB, nil).Once()
	pgDBTx.On("RollbackTx", pgDB).Once()
	pgDBStoreRepository.On("CreateStore", mock.Anything, pgDB, &mockStore).Return(nil).Once()
	pgDBQueueRepository.On("CreateQueues", mock.Anything, pgDB, 1, mockQueues).Return(nil).Once()
	pgDBTx.On("CommitTx", pgDB).Return(nil).Once()

	err := integrationUsecase.CreateStore(context.TODO(), &mockStore, mockQueues)
	assert.NoError(t, err)
}
