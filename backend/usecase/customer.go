package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (uc *Usecase) CreateCustomers(ctx context.Context, session domain.StoreSession, oldStatus string, newStatus string, customers []domain.Customer) error {
	tx, err := uc.pgDBRepository.BeginTx()
	if err != nil {
		return err
	}
	defer uc.pgDBRepository.RollbackTx(tx)

	err = uc.pgDBRepository.UpdateSessionWithTx(ctx, tx, session, oldStatus, newStatus)
	if err != nil {
		return err
	}

	err = uc.pgDBRepository.CreateCustomers(ctx, tx, customers)
	if err != nil {
		return err
	}
	
	err = uc.pgDBRepository.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error {
	err := uc.pgDBRepository.UpdateCustomer(ctx, oldStatus, newStatus, customer)
	return err
}
