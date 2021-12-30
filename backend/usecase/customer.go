package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
)

type customerUsecase struct {
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface
	logger                 logging.LoggerTool
}

func NewCustomerUsecase(
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface,
	logger logging.LoggerTool,
) CustomerUseCaseInterface {
	return &customerUsecase{pgDBCustomerRepository, logger}
}

func (cu *customerUsecase) UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error {
	err := cu.pgDBCustomerRepository.UpdateCustomer(ctx, oldStatus, newStatus, customer)
	return err
}
