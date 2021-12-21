package pgDB

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
)

type PgDBRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPGDBRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) *PgDBRepository {
	return &PgDBRepository{db, logger, contextTimeOut}
}
