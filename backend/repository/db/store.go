package repository

import (
	"database/sql"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type storeRepository struct {
	db *sql.DB
	logger logging.LoggerTool
}

func NewStoreRepository(db *sql.DB, logger logging.LoggerTool) domain.StoreRepositoryInterface {
	return &storeRepository{db, logger}
}