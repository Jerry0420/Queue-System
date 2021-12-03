package pgDB

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPGDBRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBRepositoryInterface {
	return &pgDBRepository{db, logger, contextTimeOut}
}
