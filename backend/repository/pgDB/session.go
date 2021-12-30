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
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBSessionRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBSessionRepositoryInterface {
	return &pgDBSessionRepository{db, logger, contextTimeOut}
}

func (psr *pgDBSessionRepository) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{StoreId: store.ID, StoreSessionStatus: domain.StoreSessionStatus.NORMAL}

	query := `INSERT INTO store_sessions (store_id) VALUES ($1) RETURNING id`
	stmt, err := psr.db.PrepareContext(ctx, query)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, store.ID)
	err = row.Scan(&session.ID)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError40904
	}
	return session, nil
}

func (psr *pgDBSessionRepository) UpdateSessionStatus(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	query := `UPDATE store_sessions SET status=$1 WHERE id=$2 and status=$3`
	stmt, err := psr.db.PrepareContext(ctx, query)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, newStatus, session.ID, oldStatus)
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
}

func (psr *pgDBSessionRepository) UpdateSessionWithTx(ctx context.Context, tx *sql.Tx, session domain.StoreSession, oldStatus string, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()
	query := `UPDATE store_sessions SET status=$1 WHERE id=$2 and status=$3`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, newStatus, session.ID, oldStatus)
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
}

func (psr *pgDBSessionRepository) GetSessionAndStoreBySessionId(ctx context.Context, sessionId string) (domain.StoreSession, domain.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{}
	store := domain.Store{}

	query := `SELECT stores.id, stores.created_at, store_sessions.status 
				FROM store_sessions
				INNER JOIN stores ON stores.id = store_sessions.store_id WHERE store_sessions.id=$1`
	row := psr.db.QueryRowContext(ctx, query, sessionId)
	err := row.Scan(&store.ID, &store.CreatedAt, &session.StoreSessionStatus)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		psr.logger.ERRORf("error %v", err)
		return session, store, domain.ServerError40404
	case err != nil:
		psr.logger.ERRORf("error %v", err)
		return session, store, domain.ServerError50002
	}
	session.ID = sessionId
	session.StoreId = store.ID
	return session, store, nil
}
