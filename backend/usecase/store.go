package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type storeUsecase struct {
	storeRepository domain.StoreRepositoryInterface
	logger logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, logger}
}

func (su *storeUsecase) Create(ctx context.Context, store *domain.Store) error {
	return nil
}