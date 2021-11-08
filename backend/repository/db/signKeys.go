package repository

import (
	"context"
	"database/sql"
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

func (skr *signKeyRepository) Create(ctx context.Context, signKey *domain.SignKey) (signKeyID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, skr.contextTimeOut)
	defer cancel()

	query := `INSERT INTO sign_keys (store_id, sign_key, type) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := skr.db.PrepareContext(ctx, query)
	if err != nil {
		return signKeyID, domain.ServerError50002
	}
	row := stmt.QueryRowContext(ctx, signKey.StoreId, signKey.SignKey, signKey.SignKeyType)
	err = row.Scan(&signKeyID)
	if err != nil {
		skr.logger.ERRORf("error %v", err)
		return signKeyID, domain.ServerError50002
	}
	return signKeyID, nil
}
