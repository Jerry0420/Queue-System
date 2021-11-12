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
	Update(ctx context.Context, store *Store, fieldName string, newFieldValue string) error
}

type StoreUsecaseInterface interface {
	Create(ctx context.Context, store Store) error
	GetByEmail(ctx context.Context, email string) (Store, error)
	VerifyPasswordLength(password string) error
	EncryptPassword(password string) (string, error)
	ValidatePassword(ctx context.Context, incomingPassword string, password string) error
	GenerateToken(ctx context.Context, store Store, signKeyType string, expiresDuration time.Duration) (string, error)
	VerifyToken(ctx context.Context, encryptToken string) (tokenClaims TokenClaims, err error)
	VerifyTokenRenewable(tokenClaims TokenClaims) bool
	RemoveSignKeyByID(ctx context.Context, signKeyID int) error
	GenerateEmailContentOfForgetPassword(emailToken string, store Store) (subject string, content string)
	Update(ctx context.Context, store *Store, fieldName string, newFieldValue string) error
}