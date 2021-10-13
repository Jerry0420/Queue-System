package usecase

import (
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
)

type storeUsecase struct {
	storeRepository domain.StoreRepositoryInterface
	logger logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, logger}
}