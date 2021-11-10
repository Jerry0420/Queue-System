package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
)

type storeStatus struct{ OPEN, CLOSE string }

var StoreStatus storeStatus = storeStatus{OPEN: "open", CLOSE: "close"}

type Store struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	SessionID   string    `json:"session_id"`
}

type TokenClaims struct {
	StoreID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	SignKeyID int    `json:"signkey_id"`
	jwt.StandardClaims
}

type StoreRepositoryInterface interface {
	Create(ctx context.Context, store Store) error
	GetByEmail(ctx context.Context, email string) (Store, error)
}

type StoreUsecaseInterface interface {
	Create(ctx context.Context, store Store) error
	GetByEmail(ctx context.Context, email string) (Store, error)
	Signin(ctx context.Context, store Store) (Store, error)
	GenerateToken(ctx context.Context, store Store) (string, error)
	VerifyToken(ctx context.Context, encryptToken string) (tokenClaims TokenClaims, err error)
	RemoveSignKeyByID(ctx context.Context, signKeyID int) error
}

// query := `INSERT INTO stores (email, password, name, description, status) VALUES ($1, $2, $3, $4, $5) `
// stmt, err := db.PrepareContext(r.Context(), query)
// _, err = stmt.ExecContext(r.Context(), "jeerywa@gmail.com", "im password", "jerry", "im description", "open")
// utils.JsonResponseOK(w, nil)

// query := `UPDATE stores SET name=$1,session_id=uuid_generate_v4()`
// stmt, err := db.PrepareContext(r.Context(), query)
// _, err = stmt.ExecContext(r.Context(), "new_jerry")

// query = `SELECT id,name,created_at,status,session_id FROM stores`
// rows, err := db.QueryContext(r.Context(), query)
// defer func() {
// 	errRow := rows.Close()
// }()
// result := make([]domain.Store, 0)
// for rows.Next() {
// 	store := domain.Store{}
// 	err = rows.Scan(&store.ID, &store.Name, &store.CreatedAt, &store.Status, &store.SessionID)
// 	result = append(result, store)
// }
