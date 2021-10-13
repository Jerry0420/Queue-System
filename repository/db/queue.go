package repository

import (
	"database/sql"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
)

type queueRepository struct {
	db *sql.DB
	logger logging.LoggerTool
}

func NewQueueRepository(db *sql.DB, logger logging.LoggerTool) domain.QueueRepositoryInterface {
	return &queueRepository{db, logger}
}