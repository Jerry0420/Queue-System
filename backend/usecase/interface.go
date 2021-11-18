package usecase

import (
	"context"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

type UseCaseInterface interface {
	// store.go
	CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error
	GetStoreByEmail(ctx context.Context, email string) (domain.Store, error)
	VerifyPasswordLength(password string) error
	EncryptPassword(password string) (string, error)
	ValidatePassword(ctx context.Context, incomingPassword string, password string) error
	CloseStore(ctx context.Context, store domain.Store) error
	GenerateToken(ctx context.Context, store domain.Store, signKeyType string, expireTime time.Time) (encryptToken string, err error)
	VerifyToken(
		ctx context.Context,
		encryptToken string,
		signKeyType string,
		getSignKey func(context.Context, int, string) (domain.SignKey, error),
	) (tokenClaims domain.TokenClaims, err error)
	GetSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error)
	RemoveSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error)
	GenerateEmailContentOfForgetPassword(passwordToken string, store domain.Store) (subject string, content string)
	UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error

	// queue.go

	// customer.go
}
