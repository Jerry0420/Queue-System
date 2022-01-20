package pgDB_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreByEmail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}
	mockStore := domain.Store{
		ID:          1,
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "",
		CreatedAt:   time.Now(),
		Timezone:    "Asia/Taipei",
	}
	rows := sqlmock.NewRows([]string{"id", "email", "password", "name", "description", "created_at", "timezone"}).
		AddRow(mockStore.ID, mockStore.Email, mockStore.Password, mockStore.Name, mockStore.Description, mockStore.CreatedAt, mockStore.Timezone)

	query := `SELECT id,email,password,name,description,created_at,timezone FROM stores WHERE email=$1`
	mock.ExpectQuery(query).WithArgs("email1").WillReturnRows(rows)
	logger := logging.NewLogger([]string{}, true)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, time.Duration(2*time.Second))

	t.Run("right email", func(t *testing.T) {
		store, err := pgDBStoreRepository.GetStoreByEmail(context.TODO(), "email1")
		assert.NoError(t, err)
		assert.Equal(t, mockStore, store)
	})

	t.Run("wrong email", func(t *testing.T) {
		store, err := pgDBStoreRepository.GetStoreByEmail(context.TODO(), "email2")
		assert.NotNil(t, err)
		assert.Equal(t, domain.Store{}, store)
	})

}

func TestCreateStore(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}
	mockStore := domain.Store{
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "",
		Timezone:    "Asia/Taipei",
	}

	query := `INSERT INTO stores (name, email, password, timezone) VALUES ($1, $2, $3, $4) RETURNING id,created_at`

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs(mockStore.Name, mockStore.Email, mockStore.Password, mockStore.Timezone).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(66, time.Now()))

	logger := logging.NewLogger([]string{}, true)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, time.Duration(2*time.Second))
	err = pgDBStoreRepository.CreateStore(context.TODO(), db, &mockStore)
	assert.NoError(t, err)
	assert.Equal(t, 66, mockStore.ID)
}
