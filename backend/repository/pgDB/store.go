package pgDB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (repo *PgDBRepository) GetStoreByEmail(ctx context.Context, email string) (domain.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	query := `SELECT id,email,password,name,description,created_at FROM stores WHERE email=$1`
	row := repo.db.QueryRowContext(ctx, query, email)
	var store domain.Store
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

func (repo *PgDBRepository) GetStoreWIthQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	var storeWithQueues domain.StoreWithQueues
	query := `SELECT 
					stores.email, stores.name, stores.description, stores.created_at, 
					queues.id AS queue_id, queues.name AS queue_name, 
					customers.id AS customer_id, customers.name AS customer_name, customers.phone AS customer_phone, 
					customers.status AS customer_status,
					customers.created_at AS customer_created_at
			FROM stores
			INNER JOIN queues ON stores.id = queues.store_id
			INNER JOIN customers ON queues.id = customers.queue_id
			WHERE stores.id=$1 and customers.status='normal' OR customers.status='processing'
			ORDER BY customers.id ASC`

	rows, err := repo.db.QueryContext(ctx, query, storeId)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return storeWithQueues, domain.ServerError50002
	}

	var store domain.Store
	queues := make(map[int]domain.Queue)
	customers := make(map[int][]domain.Customer)
	
	for rows.Next() {
		var queue domain.Queue
		var customer domain.Customer

		err := rows.Scan(
			&store.Email, &store.Name, &store.Description, &store.CreatedAt,
			&queue.ID, &queue.Name,
			&customer.ID, &customer.Name, &customer.Phone, &customer.Status, &customer.CreatedAt,
		)
		if err != nil {
			repo.logger.ERRORf("error %v", err)
			return storeWithQueues, domain.ServerError50002
		}
		queues[queue.ID] = queue
		customers[queue.ID] = append(customers[queue.ID], customer)
	}
	defer rows.Close()

	if store == (domain.Store{}) {
		queues = make(map[int]domain.Queue)
		query = `SELECT 
					stores.email, 
					stores.name, 
					stores.description, 
					stores.created_at, 
					queues.id AS queue_id, 
					queues.name AS queue_name
				FROM stores
				INNER JOIN queues ON stores.id = queues.store_id
				WHERE stores.id=$1`
		rows, err = repo.db.QueryContext(ctx, query, storeId)
		if err != nil {
			repo.logger.ERRORf("error %v", err)
			return storeWithQueues, domain.ServerError50002
		}
		for rows.Next() {
			var queue domain.Queue
			err := rows.Scan(
				&store.Email, &store.Name, &store.Description, &store.CreatedAt,
				&queue.ID, &queue.Name,
			)
			if err != nil {
				repo.logger.ERRORf("error %v", err)
				return storeWithQueues, domain.ServerError50002
			}
			queues[queue.ID] = queue
		}
	}

	if store == (domain.Store{}) {
		return storeWithQueues, domain.ServerError40402
	}

	storeWithQueues = domain.StoreWithQueues{ID: storeId, Email: store.Email, Name: store.Name, Description: store.Description, CreatedAt: store.CreatedAt}
	for _, queue := range queues {
		storeWithQueues.Queues = append(storeWithQueues.Queues, domain.QueueWithCustomers{
			ID: queue.ID,
			Name: queue.Name,
			Customers: customers[queue.ID],
		})
	}
	return storeWithQueues, nil
}

func (repo *PgDBRepository) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
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

func (repo *PgDBRepository) UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
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

func (repo *PgDBRepository) RemoveStoreByID(ctx context.Context, id int) error {
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
