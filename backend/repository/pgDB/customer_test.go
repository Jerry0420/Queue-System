package pgDB_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/stretchr/testify/assert"
)

func setUpCustomerTest(t *testing.T) (pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface, db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}

	logger := logging.NewLogger([]string{}, true)
	pgDBCustomerRepository = pgDB.NewPgDBCustomerRepository(db, logger, time.Duration(2*time.Second))
	return pgDBCustomerRepository, db, mock
}

func TestUpdateCustomer(t *testing.T) {
	pgDBCustomerRepository, _, mock := setUpCustomerTest(t)
	oldMockCustomerStatus := domain.CustomerStatus.PROCESSING
	newMockCustomerStatus := domain.CustomerStatus.DONE
	mockCustomer := domain.Customer{
		ID: 1,
		Name: "customer1",
		Phone: "0000000000",
		QueueID: 1,
		Status: oldMockCustomerStatus,
		CreatedAt: time.Now(),
	}
	query := `UPDATE customers SET status=$1 WHERE id=$2 and status=$3`
	mock.ExpectExec(query).
		WithArgs(newMockCustomerStatus, mockCustomer.ID, oldMockCustomerStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := pgDBCustomerRepository.UpdateCustomer(context.TODO(), oldMockCustomerStatus, newMockCustomerStatus, &mockCustomer)
	assert.NoError(t, err)
}
