package pgDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

// interface for sql.DB and sql.Tx
type PgDBInterface interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type PgDBTxInterface interface {
	BeginTx() (tx PgDBInterface, err error)
	RollbackTx(pgDbTx PgDBInterface)
	CommitTx(pgDbTx PgDBInterface) error
}

type PgDBStoreRepositoryInterface interface {
	GetStoreByEmail(ctx context.Context, email string) (domain.Store, error)
	GetStoreWithQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error)
	GetStoreWithQueuesById(ctx context.Context, storeId int) (domain.StoreWithQueues, error)
	CreateStore(ctx context.Context, tx PgDBInterface, store *domain.Store, queues []domain.Queue) error
	UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error
	RemoveStoreByID(ctx context.Context, tx PgDBInterface, id int) error
	RemoveStoreByIDs(ctx context.Context, tx PgDBInterface, storeIds []string) error
	GetAllIdsOfExpiredStores(ctx context.Context, tx PgDBInterface, expiresTime time.Time) (storesIds []string, err error)
	GetAllExpiredStoresInSlice(ctx context.Context, tx PgDBInterface, expiresTime time.Time) (stores [][][]string, err error)
}

type PgDBSignKeyRepositoryInterface interface {
	CreateSignKey(ctx context.Context, signKey *domain.SignKey) error
	GetSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)
	RemoveSignKeyByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error)
}

type PgDBQueueRepositoryInterface interface {
	CreateQueues(ctx context.Context, tx PgDBInterface, storeID int, queues []domain.Queue) error
}

type PgDBCustomerRepositoryInterface interface {
	CreateCustomers(ctx context.Context, tx PgDBInterface, customers []domain.Customer) error
	UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error
	GetCustomersWithQueuesByStoreId(ctx context.Context, tx PgDBInterface, storeId int) (customers [][]string, err error)
}

type PgDBSessionRepositoryInterface interface {
	CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error)
	UpdateSessionStatus(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error
	UpdateSessionWithTx(ctx context.Context, tx PgDBInterface, session domain.StoreSession, oldStatus string, newStatus string) error
	GetSessionAndStoreBySessionId(ctx context.Context, sessionId string) (domain.StoreSession, domain.Store, error)
}
