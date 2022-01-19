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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}
	mockStore := domain.Store{
		ID:          1,
		Email:       "email 1",
		Password:    "password 1",
		Name:        "name 1",
		Description: "",
		CreatedAt:   time.Now(),
		Timezone:    "Asia/Taipei",
	}
	rows := sqlmock.NewRows([]string{"id", "email", "password", "name", "description", "created_at", "timezone"}).
		AddRow(mockStore.ID, mockStore.Email, mockStore.Password, mockStore.Name, mockStore.Description, mockStore.CreatedAt, mockStore.Timezone)

	query := `SELECT id,email,password,name,description,created_at,timezone FROM stores WHERE email=\\?`
	mock.ExpectQuery(query).WithArgs("email 1").WillReturnRows(rows)
	logger := logging.NewLogger([]string{}, true)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, time.Duration(2*time.Second))

	t.Run("right email", func(t *testing.T) {
		store, err := pgDBStoreRepository.GetStoreByEmail(context.TODO(), "email 1")
		assert.NoError(t, err)
		assert.Equal(t, mockStore, store)
	})

	t.Run("wrong email", func(t *testing.T) {
		store, err := pgDBStoreRepository.GetStoreByEmail(context.TODO(), "email 2")
		assert.NotNil(t, err)
		assert.Equal(t, domain.Store{}, store)
	})

}
