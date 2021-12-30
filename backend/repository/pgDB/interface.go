package pgDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

type PgDBStoreRepositoryInterface interface {
	GetStoreByEmail(ctx context.Context, email string) (domain.Store, error)
	GetStoreWithQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error)
	GetStoreWithQueuesById(ctx context.Context, storeId int) (domain.StoreWithQueues, error)
	CreateStore(ctx context.Context, tx *sql.Tx, store *domain.Store, queues []domain.Queue) error
	UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error
	RemoveStoreByID(ctx context.Context, tx *sql.Tx, id int) error
	RemoveStoreByIDs(ctx context.Context, tx *sql.Tx, storeIds []string) error
	GetAllIdsOfExpiredStores(ctx context.Context, tx *sql.Tx, expiresTime time.Time) (storesIds []string, err error)
	GetAllExpiredStoresInSlice(ctx context.Context, tx *sql.Tx, expiresTime time.Time) (stores[][][]string, err error)
}

type PgDBSignKeyRepositoryInterface interface {
	CreateSignKey(ctx context.Context, signKey *domain.SignKey) error
	GetSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)
	RemoveSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)
}

type PgDBQueueRepositoryInterface interface {
	CreateQueues(ctx context.Context, tx *sql.Tx, storeID int, queues []domain.Queue) error
}

type PgDBCustomerRepositoryInterface interface {
	CreateCustomers(ctx context.Context, tx *sql.Tx, customers []domain.Customer) error
	UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error
	GetCustomersWithQueuesByStoreId(ctx context.Context, tx *sql.Tx, storeId int) (customers [][]string, err error)
}

type PgDBSessionRepositoryInterface interface {
	CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error)
	UpdateSessionStatus(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error
	UpdateSessionWithTx(ctx context.Context, tx *sql.Tx, session domain.StoreSession, oldStatus string, newStatus string) error
	GetSessionAndStoreBySessionId(ctx context.Context, sessionId string) (domain.StoreSession, domain.Store, error)
}
