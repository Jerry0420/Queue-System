package pgDB

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (repo *pgDBRepository) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{StoreId: store.ID, StoreSessionStatus: domain.StoreSessionStatus.NORMAL}

	query := `INSERT INTO store_sessions (store_id) VALUES ($1) RETURNING id`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return session, domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, store.ID)
	err = row.Scan(&session.ID)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return session, domain.ServerError40904
	}
	return session, nil
}

func (repo *pgDBRepository) UpdateSession(ctx context.Context, session domain.StoreSession, oldStatus string, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	query := `UPDATE store_sessions SET status=$1 WHERE id=$2 and status=$3`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, newStatus, session.ID, oldStatus)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40404
	}
	return nil
}

func (repo *pgDBRepository) UpdateSessionWithTx(ctx context.Context, tx *sql.Tx, session domain.StoreSession, oldStatus string, newStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()
	query := `UPDATE store_sessions SET status=$1 WHERE id=$2 and status=$3`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, newStatus, session.ID, oldStatus)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40404
	}
	return nil
}

func (repo *pgDBRepository) GetSessionById(ctx context.Context, sessionId string) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{ID: sessionId}
	query := `SELECT store_id FROM store_sessions WHERE id=$1`
	row := repo.db.QueryRowContext(ctx, query, sessionId)
	err := row.Scan(&session.StoreId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		repo.logger.ERRORf("error %v", err)
		return session, domain.ServerError40404
	case err != nil:
		repo.logger.ERRORf("error %v", err)
		return session, domain.ServerError50002
	}
	return session, nil
}
