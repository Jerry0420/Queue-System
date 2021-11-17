package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
)

type Store struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	SessionID   string    `json:"session_id"`
}

type TokenClaims struct {
	StoreID        int    `json:"store_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	StoreCreatedAt int64  `json:"store_created_at"`
	SignKeyID      int    `json:"signkey_id"`
	jwt.StandardClaims
}

type StoreRepositoryInterface interface {
	Create(ctx context.Context, store *Store, queues []Queue) error
	GetByEmail(ctx context.Context, email string) (Store, error)
	Update(ctx context.Context, store *Store, fieldName string, newFieldValue string) error
	RemoveByID(ctx context.Context, id int) error
}

type StoreUsecaseInterface interface {
	Create(ctx context.Context, store *Store, queues []Queue) error
	GetByEmail(ctx context.Context, email string) (Store, error)
	VerifyPasswordLength(password string) error
	EncryptPassword(password string) (string, error)
	ValidatePassword(ctx context.Context, incomingPassword string, password string) error
	Close(ctx context.Context, store Store) error
	GenerateToken(ctx context.Context, store Store, signKeyType string, expireTime time.Time) (encryptToken string, err error)
	VerifyToken(
		ctx context.Context,
		encryptToken string,
		signKeyType string,
		getSignKey func(context.Context, int, string) (SignKey, error),
	) (tokenClaims TokenClaims, err error)
	GetSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (SignKey, error)
	RemoveSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (SignKey, error)
	GenerateEmailContentOfForgetPassword(passwordToken string, store Store) (subject string, content string)
	Update(ctx context.Context, store *Store, fieldName string, newFieldValue string) error
}
