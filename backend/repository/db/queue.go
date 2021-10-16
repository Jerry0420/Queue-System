package repository

import (
	"database/sql"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type queueRepository struct {
	db *sql.DB
	logger logging.LoggerTool
}

func NewQueueRepository(db *sql.DB, logger logging.LoggerTool) domain.QueueRepositoryInterface {
	return &queueRepository{db, logger}
}