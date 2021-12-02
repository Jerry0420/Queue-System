package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (uc *usecase) CreateCustomer(ctx context.Context, session domain.StoreSession, oldStatus string, newStatus string, customers []domain.Customer) error {
	err := uc.pgDBRepository.CreateCustomers(ctx, session, oldStatus, newStatus, customers)
	return err
}
