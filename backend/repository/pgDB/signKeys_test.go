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

func setUpSignKeysTest(t *testing.T) (pgDBSignKeyRepository pgDB.PgDBSignKeyRepositoryInterface, db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}

	logger := logging.NewLogger([]string{}, true)
	pgDBSignKeyRepository = pgDB.NewPgDBSignKeyRepository(db, logger, time.Duration(2*time.Second))
	return pgDBSignKeyRepository, db, mock
}

func TestCreateSignKey(t *testing.T) {
	pgDBSignKeyRepository, _, mock := setUpSignKeysTest(t)
	mockSingKey := domain.SignKey{
		StoreId: 1,
		SignKey: "imsignkey",
		SignKeyType: domain.SignKeyTypes.NORMAL,
	}

	query := `INSERT INTO sign_keys (store_id, sign_key, type) VALUES ($1, $2, $3) RETURNING id`
	mock.ExpectQuery(query).
		WithArgs(mockSingKey.StoreId, mockSingKey.SignKey, mockSingKey.SignKeyType).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(66))

	err := pgDBSignKeyRepository.CreateSignKey(context.TODO(), &mockSingKey)
	assert.NoError(t, err)
	assert.Equal(t, 66, mockSingKey.ID)
}
