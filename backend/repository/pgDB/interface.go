package pgDB

import (
	"context"
	"database/sql"

	"github.com/jerry0420/queue-system/backend/domain"
)

type PgDBRepositoryInterface interface {
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
	CreateQueues(ctx context.Context, tx *sql.Tx, storeID int, queues []domain.Queue) error

	// customer.go

	// session.go
	CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error)
	UpdateSession(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error
}
