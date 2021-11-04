package repository

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type queueRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewQueueRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) domain.QueueRepositoryInterface {
	return &queueRepository{db, logger, contextTimeOut}
}
