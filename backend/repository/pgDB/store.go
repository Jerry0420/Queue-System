package pgDB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (repo *pgDBRepository) GetStoreByEmail(ctx context.Context, email string) (domain.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	var store domain.Store
	query := `SELECT id,email,password,name,description,created_at FROM stores WHERE email=$1`
	row := repo.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&store.ID, &store.Email, &store.Password, &store.Name, &store.Description, &store.CreatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		repo.logger.ERRORf("error %v", err)
		return store, domain.ServerError40402
	case err != nil:
		repo.logger.ERRORf("error %v", err)
		return store, domain.ServerError50002
	}
	return store, nil
}

func (repo *pgDBRepository) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO stores (name, email, password) VALUES ($1, $2, $3) RETURNING id,created_at`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, store.Name, store.Email, store.Password)
	err = row.Scan(&store.ID, &store.CreatedAt)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError40901
	}

	err = repo.CreateQueues(ctx, tx, store.ID, queues)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	return nil
}

func (repo *pgDBRepository) UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	query := fmt.Sprintf("UPDATE stores SET %s=$1 WHERE id=$2 RETURNING description,created_at", fieldName)
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, newFieldValue, store.ID)
	err = row.Scan(&store.Description, &store.CreatedAt)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError40402
	}
	return nil
}

func (repo *pgDBRepository) RemoveStoreByID(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	query := `DELETE FROM stores WHERE id=$1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
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
		return domain.ServerError40402
	}
	return nil
}
