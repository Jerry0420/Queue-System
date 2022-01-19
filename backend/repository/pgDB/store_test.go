package pgDB_test

import (
	"context"
	"fmt"
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
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "",
		CreatedAt:   time.Now(),
		Timezone:    "Asia/Taipei",
	}
	rows := sqlmock.NewRows([]string{"id", "email", "password", "name", "description", "created_at", "timezone"}).
		AddRow(mockStore.ID, mockStore.Email, mockStore.Password, mockStore.Name, mockStore.Description, mockStore.CreatedAt, mockStore.Timezone)

	query := `SELECT id,email,password,name,description,created_at,timezone FROM stores WHERE email=\\?`
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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}
	mockStore := domain.Store{
		Email:       "email 1",
		Password:    "password 1",
		Name:        "name 1",
		Description: "",
		Timezone:    "Asia/Taipei",
	}
	
	query := `INSERT INTO stores (name, email, password, timezone) VALUES (\\?, \\?, \\?, \\?) RETURNING id,created_at`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(mockStore.Name, mockStore.Email, mockStore.Password, mockStore.Timezone).WillReturnResult(sqlmock.NewResult(1, 1))

	logger := logging.NewLogger([]string{}, true)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, time.Duration(2*time.Second))
	err = pgDBStoreRepository.CreateStore(context.TODO(), db, &mockStore)
	fmt.Println(`1=======================================================`)
	fmt.Println(err)
	fmt.Println(`2=======================================================`)
	fmt.Println(mockStore)
	fmt.Println(`3=======================================================`)
	assert.Equal(t, 1, 1)
}
