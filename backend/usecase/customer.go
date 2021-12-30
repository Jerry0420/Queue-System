package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
)

type customerUsecase struct {
	pgDBTx                 pgDB.PgDBTxInterface
	pgDBSessionRepository  pgDB.PgDBSessionRepositoryInterface
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface
	logger                 logging.LoggerTool
}

func NewCustomerUsecase(
	pgDBTx pgDB.PgDBTxInterface,
	pgDBSessionRepository pgDB.PgDBSessionRepositoryInterface,
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface,
	logger logging.LoggerTool,
) CustomerUseCaseInterface {
	return &customerUsecase{pgDBTx, pgDBSessionRepository, pgDBCustomerRepository, logger}
}

func (cu *customerUsecase) CreateCustomers(ctx context.Context, session domain.StoreSession, oldStatus string, newStatus string, customers []domain.Customer) error {
	tx, err := cu.pgDBTx.BeginTx()
	if err != nil {
		return err
	}
	defer cu.pgDBTx.RollbackTx(tx)

	err = cu.pgDBSessionRepository.UpdateSessionWithTx(ctx, tx, session, oldStatus, newStatus)
	if err != nil {
		return err
	}

	err = cu.pgDBCustomerRepository.CreateCustomers(ctx, tx, customers)
	if err != nil {
		return err
	}

	err = cu.pgDBTx.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (cu *customerUsecase) UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error {
	err := cu.pgDBCustomerRepository.UpdateCustomer(ctx, oldStatus, newStatus, customer)
	return err
}
