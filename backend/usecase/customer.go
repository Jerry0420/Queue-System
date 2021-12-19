package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (uc *usecase) CreateCustomers(ctx context.Context, session domain.StoreSession, oldStatus string, newStatus string, customers []domain.Customer) error {
	err := uc.pgDBRepository.CreateCustomers(ctx, session, oldStatus, newStatus, customers)
	return err
}

func (uc *usecase) UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error {
	err := uc.pgDBRepository.UpdateCustomer(ctx, oldStatus, newStatus, customer)
	return err
}
