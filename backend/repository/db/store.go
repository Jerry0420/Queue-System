package repository

import (
	"context"
	"database/sql"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type storeRepository struct {
	db *sql.DB
	logger logging.LoggerTool
}

func NewStoreRepository(db *sql.DB, logger logging.LoggerTool) domain.StoreRepositoryInterface {
	return &storeRepository{db, logger}
}

func (sr *storeRepository) Create(ctx context.Context, store *domain.Store) error {
	return nil
}