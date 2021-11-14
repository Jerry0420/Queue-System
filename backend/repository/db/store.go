package repository

import (
	"context"
	"database/sql"
	"errors"
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
	query := `SELECT id,email,password,name,description,created_at,session_id FROM stores WHERE email=$1`
	row := sr.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&store.ID, &store.Email, &store.Password, &store.Name, &store.Description, &store.CreatedAt, &store.SessionID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		sr.logger.ERRORf("error %v", err)
		return store, domain.ServerError40402
	case err != nil:
		sr.logger.ERRORf("error %v", err)
		return store, domain.ServerError50002
	}
	return store, nil
}

func (sr *storeRepository) Create(ctx context.Context, store *domain.Store) error {
	ctx, cancel := context.WithTimeout(ctx, sr.contextTimeOut)
	defer cancel()
	
	query := `INSERT INTO stores (name, email, password) VALUES ($1, $2, $3) RETURNING id,created_at,session_id`
	stmt, err := sr.db.PrepareContext(ctx, query)
	if err != nil {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	row := stmt.QueryRowContext(ctx, store.Name, store.Email, store.Password)
	err = row.Scan(&store.ID, &store.CreatedAt, &store.SessionID)
	if err != nil {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError40402
	}
	return nil
}

func (sr *storeRepository) Update(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
	ctx, cancel := context.WithTimeout(ctx, sr.contextTimeOut)
	defer cancel()

	query := fmt.Sprintf("UPDATE stores SET %s=$1 WHERE id=$2 RETURNING description,created_at,session_id", fieldName)
	stmt, err := sr.db.PrepareContext(ctx, query)
	if err != nil {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	row := stmt.QueryRowContext(ctx, newFieldValue, store.ID)
	err = row.Scan(&store.Description, &store.CreatedAt, &store.SessionID)
	if err != nil {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError40402
	}
	return nil
}

func (sr *storeRepository) RemoveByID(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, sr.contextTimeOut)
	defer cancel()

	var deletedID int
	query := `DELETE FROM stores WHERE id=$1 RETURNING id`
	stmt, err := sr.db.PrepareContext(ctx, query)
	if err != nil {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&deletedID)
	if err != nil || deletedID != id {
		sr.logger.ERRORf("error %v", err)
		return domain.ServerError40403
	}
	return nil
}