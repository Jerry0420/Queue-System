package repository

import (
	"database/sql"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
)

type customerRepository struct {
	db *sql.DB
	logger logging.LoggerTool
}

func NewCustomerRepository(db *sql.DB, logger logging.LoggerTool) domain.CustomerRepositoryInterface {
	return &customerRepository{db, logger}
}