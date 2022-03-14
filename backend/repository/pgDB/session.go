package pgDB

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBSessionRepository struct {
	db             PgDBInterface
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBSessionRepository(db PgDBInterface, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBSessionRepositoryInterface {
	return &pgDBSessionRepository{db, logger, contextTimeOut}
}

func (psr *pgDBSessionRepository) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{StoreId: store.ID, StoreSessionStatus: domain.StoreSessionStatus.NORMAL}

	query := `INSERT INTO store_sessions (store_id) VALUES ($1) RETURNING id`
	row := psr.db.QueryRowContext(ctx, query, store.ID)
	err := row.Scan(&session.ID)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError40903
	}
	return session, nil
}

func (psr *pgDBSessionRepository) UpdateSessionStatus(ctx context.Context, tx PgDBInterface, session *domain.StoreSession, oldStatus string, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	sessionStatusInDb := ""
	var err error
	var row *sql.Row

	query := `SELECT status FROM store_sessions WHERE id=$1`
	if tx == nil {
		row = psr.db.QueryRowContext(ctx, query, session.ID)
	} else {
		row = tx.QueryRowContext(ctx, query, session.ID)
	}
	err = row.Scan(&sessionStatusInDb)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError40404
	case err != nil:
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}

	if sessionStatusInDb == oldStatus {
		query = `UPDATE store_sessions SET status=$1 WHERE id=$2 and status=$3`
		var result sql.Result
		if tx == nil {
			result, err = psr.db.ExecContext(ctx, query, newStatus, session.ID, oldStatus)
		} else {
			result, err = tx.ExecContext(ctx, query, newStatus, session.ID, oldStatus)
		}
		if err != nil {
			psr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		num, err := result.RowsAffected()
		if err != nil {
			psr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		if num == 0 {
			return domain.ServerError40404
		}
		return nil
	} else {
		switch sessionStatusInDb {
			case domain.StoreSessionStatus.SCANNED: 
				return domain.ServerError40007
			case domain.StoreSessionStatus.USED: 
				return domain.ServerError40008
		}
	}

	return nil
}

func (psr *pgDBSessionRepository) GetSessionById(ctx context.Context, sessionId string) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{}
	store := domain.Store{}

	query := `SELECT stores.id, store_sessions.status 
				FROM store_sessions
				INNER JOIN stores ON stores.id = store_sessions.store_id WHERE store_sessions.id=$1`
	row := psr.db.QueryRowContext(ctx, query, sessionId)
	err := row.Scan(&store.ID, &session.StoreSessionStatus)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError40404
	case err != nil:
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError50002
	}
	session.ID = sessionId
	session.StoreId = store.ID
	return session, nil
}
