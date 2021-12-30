package pgDB

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBCustomerRepository struct {
	db             *sql.DB
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBCustomerRepository(db *sql.DB, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBCustomerRepositoryInterface {
	return &pgDBCustomerRepository{db, logger, contextTimeOut}
}

func (pcr *pgDBCustomerRepository) CreateCustomers(ctx context.Context, tx *sql.Tx, customers []domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	variableCounts := 1
	var query bytes.Buffer
	var queryRowParams []interface{}
	query.WriteString("INSERT INTO customers (name, phone, queue_id, status) VALUES ")
	for index, customer := range customers {
		query.WriteString("($")
		query.WriteString(strconv.Itoa(variableCounts))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 1))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 2))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 3))
		query.WriteString(")")
		variableCounts = variableCounts + 4
		queryRowParams = append(queryRowParams, customer.Name, customer.Phone, customer.QueueID, customer.Status)
		if index != len(customers)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" RETURNING id,name,phone,queue_id,created_at")

	stmt, err := tx.PrepareContext(ctx, query.String())
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, queryRowParams...)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	customers = customers[:0] // clear customers slice

	for rows.Next() {
		customer := domain.Customer{}
		err = rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.QueueID, &customer.CreatedAt)
		if err != nil {
			pcr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		customer.Status = domain.CustomerStatus.NORMAL
		customers = append(customers, customer)

	}
	defer rows.Close()

	return nil
}

func (pcr *pgDBCustomerRepository) UpdateCustomer(ctx context.Context, oldStatus string, newStatus string, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	query := `UPDATE customers SET status=$1 WHERE id=$2 and status=$3`
	stmt, err := pcr.db.PrepareContext(ctx, query)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, newStatus, customer.ID, oldStatus)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40405
	}
	return nil
}

func (pcr *pgDBCustomerRepository) GetCustomersWithQueuesByStoreId(ctx context.Context, tx *sql.Tx, storeId int) (customers [][]string, err error) {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	customers = make([][]string, 0)

	query := `SELECT 
					queues.name AS queue_name, 
					customers.name AS customer_name, customers.phone AS customer_phone,
					customers.status AS customer_status, customers.created_at AS customer_created_at
				FROM queues
				INNER JOIN customers ON queues.id = customers.queue_id
				WHERE queues.store_id=$1
				ORDER BY queues.id ASC, customers.id ASC FOR UPDATE`

	rows, err := pcr.db.QueryContext(ctx, query, storeId)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return customers, domain.ServerError50002
	}

	for rows.Next() {
		var queue domain.Queue
		var customer domain.Customer
		err := rows.Scan(
			&queue.Name,
			&customer.Name, &customer.Phone, &customer.Status, &customer.CreatedAt,
		)
		if err != nil {
			pcr.logger.ERRORf("error %v", err)
			return customers, domain.ServerError50002
		}
		if len(customers) == 0 {
			customers = [][]string{
				[]string{
					"queue_name",
					"customer_name",
					"customer_phone",
					"customer_status",
					"customer_created_at",
				},
				[]string{
					queue.Name,
					customer.Name,
					customer.Phone,
					customer.Status,
					customer.CreatedAt.Local().String(),
				},
			}
		} else {
			customers = append(customers, []string{
				queue.Name,
				customer.Name,
				customer.Phone,
				customer.Status,
				customer.CreatedAt.Local().String(),
			})
		}
	}
	defer rows.Close()
	return customers, nil
}
