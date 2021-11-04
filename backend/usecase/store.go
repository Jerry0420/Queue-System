package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"golang.org/x/crypto/bcrypt"
)

type storeUsecase struct {
	storeRepository domain.StoreRepositoryInterface
	logger          logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, logger}
}

func (su *storeUsecase) GetByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, serverError := su.storeRepository.GetByEmail(ctx, email)
	return store, serverError
}

func (su *storeUsecase) Create(ctx context.Context, store *domain.Store) error {
	storeFromDb, err := su.GetByEmail(ctx, store.Email)
	if err != nil {
		return err
	}
	if storeFromDb != (domain.Store{}) {
		return domain.ServerError40901
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(store.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ServerError50001
	}
	store.Password = string(cryptedPassword)
	store.Status = domain.StoreStatus.OPEN
	err = su.storeRepository.Create(ctx, store)
	return err
}
