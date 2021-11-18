package pgDB

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (repo *pgDBRepository) CreateQueues(ctx context.Context, tx *sql.Tx, storeID int, queues []domain.Queue) error {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	variableCounts := 1
	var query bytes.Buffer
	var queryRowParams []interface{}
	query.WriteString("INSERT INTO queues (name, store_id) VALUES ")
	for index, queue := range queues {
		query.WriteString("($")
		query.WriteString(strconv.Itoa(variableCounts))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 1))
		query.WriteString(")")
		variableCounts = variableCounts + 2
		queryRowParams = append(queryRowParams, queue.Name, storeID)
		if index != len(queues)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" RETURNING id,name")

	stmt, err := tx.PrepareContext(ctx, query.String())
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, queryRowParams...)
	queues = queues[:0] // clear queues slice

	for rows.Next() {
		queue := domain.Queue{}
		err = rows.Scan(&queue.ID, &queue.Name)
		if err != nil {
			repo.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		queue.StoreID = storeID
		queues = append(queues, queue)

	}
	defer rows.Close()

	return nil
}
