package pgDB

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository"
)

type pgDBRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPGDBRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) repository.RepositoryInterface {
	return &pgDBRepository{db, logger, contextTimeOut}
}