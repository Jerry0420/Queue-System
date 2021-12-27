package pgDB

import (
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
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
