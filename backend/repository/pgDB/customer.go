package pgDB

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type customerRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewCustomerRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) domain.CustomerRepositoryInterface {
	return &customerRepository{db, logger, contextTimeOut}
}
