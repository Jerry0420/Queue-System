package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type storeRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewStoreRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) domain.StoreRepositoryInterface {
	return &storeRepository{db, logger, contextTimeOut}
}

func (sr *storeRepository) GetByEmail(ctx context.Context, email string) (domain.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, sr.contextTimeOut)
	defer cancel()

	var store domain.Store
	query := `SELECT id,email,password,name,description,created_at,status,session_id FROM stores WHERE email=$1`
	row := sr.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&store.ID, &store.Email, &store.Password, &store.Name, &store.Description, &store.CreatedAt, &store.Status, &store.SessionID)
	switch {
	case err == sql.ErrNoRows:
		return store, nil
	case err != nil:
		return store, domain.ServerError50002
	}
	return store, nil
}

func (sr *storeRepository) Create(ctx context.Context, store *domain.Store) error {
	ctx, cancel := context.WithTimeout(ctx, sr.contextTimeOut)
	defer cancel()

	query := `INSERT INTO stores (name, email, password, status) VALUES ($1, $2, $3, $4)`
	stmt, err := sr.db.PrepareContext(ctx, query)
	if err != nil {
		return domain.ServerError50002
	}
	_, err = stmt.ExecContext(ctx, store.Name, store.Email, store.Password, store.Status)
	if err != nil {
		return domain.ServerError50002
	}
	return nil
}
