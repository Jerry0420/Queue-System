package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type signKeyRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewSignKeyRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) domain.SignKeyRepositoryInterface {
	return &signKeyRepository{db, logger, contextTimeOut}
}

func (skr *signKeyRepository) Create(ctx context.Context, signKey *domain.SignKey) error {
	ctx, cancel := context.WithTimeout(ctx, skr.contextTimeOut)
	defer cancel()

	query := `INSERT INTO sign_keys (store_id, sign_key, type) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := skr.db.PrepareContext(ctx, query)
	if err != nil {
		skr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, signKey.StoreId, signKey.SignKey, signKey.SignKeyType)
	err = row.Scan(&signKey.ID)
	if err != nil {
		skr.logger.ERRORf("error %v", err)
		return domain.ServerError40902
	}
	return nil
}

func (skr *signKeyRepository) GetByID(ctx context.Context, id int, signKeyType string) (signKey domain.SignKey, err error) {
	ctx, cancel := context.WithTimeout(ctx, skr.contextTimeOut)
	defer cancel()

	query := `SELECT sign_key FROM sign_keys WHERE id=$1,type=$2`
	row := skr.db.QueryRowContext(ctx, query, id, signKeyType)
	err = row.Scan(&signKey.SignKey)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		skr.logger.ERRORf("error %v", err)
		return signKey, domain.ServerError40403
	case err != nil:
		skr.logger.ERRORf("error %v", err)
		return signKey, domain.ServerError50002
	}
	return signKey, nil
}

func (skr *signKeyRepository) RemoveByID(ctx context.Context, id int, signKeyType string)  (signKey domain.SignKey, err error) {
	ctx, cancel := context.WithTimeout(ctx, skr.contextTimeOut)
	defer cancel()

	query := `DELETE FROM sign_keys WHERE id=$1,type=$2 RETURNING id,sign_key`
	stmt, err := skr.db.PrepareContext(ctx, query)
	if err != nil {
		skr.logger.ERRORf("error %v", err)
		return signKey, domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id, signKeyType)
	err = row.Scan(&signKey.ID, &signKey.SignKey)
	if err != nil || signKey.ID != id {
		skr.logger.ERRORf("error %v", err)
		return signKey, domain.ServerError40403
	}
	return signKey, nil
}
