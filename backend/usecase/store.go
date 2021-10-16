package usecase

import (
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type storeUsecase struct {
	storeRepository domain.StoreRepositoryInterface
	logger logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, logger}
}