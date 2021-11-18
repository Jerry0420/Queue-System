package repository

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
)

type RepositoryInterface interface {
	// store.go
	CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error
	GetStoreByEmail(ctx context.Context, email string) (domain.Store, error)
	UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error
	RemoveStoreByID(ctx context.Context, id int) error

	// signKey.go
	CreateSignKey(ctx context.Context, signKey *domain.SignKey) error
	GetSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)
	RemoveSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)

	// queue.go

	// customer.go
}