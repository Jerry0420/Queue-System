package pgDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

// interface for sql.DB and sql.Tx
type PgDbHandleInterface interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type PgDbHandle struct {
	db *sql.DB
}

type PgDBRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPGDBRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) *PgDBRepository {
	return &PgDBRepository{db, logger, contextTimeOut}
}

func (repo *PgDBRepository) BeginTx() (tx *sql.Tx, err error) {
	// without ctx
	tx, err = repo.db.Begin()
	if err != nil {
		repo.logger.ERRORf("begin tx error %v", err)
		return nil, domain.ServerError50002
	}
	return tx, nil
}

func (repo *PgDBRepository) RollbackTx(tx *sql.Tx) {
	_ = tx.Rollback()
}

func (repo *PgDBRepository) CommitTx(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		repo.logger.ERRORf("commit error %v", err)
		return domain.ServerError50002
	}
	return nil
}
