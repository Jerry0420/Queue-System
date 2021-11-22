package pgDB

import (
	"context"

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
