package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type storeUsecase struct {
	storeRepository domain.StoreRepositoryInterface
	logger          logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, logger}
}

func (su *storeUsecase) GetByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, serverError := su.storeRepository.GetByEmail(ctx, email)
	return store, serverError
}

func (su *storeUsecase) Create(ctx context.Context, store *domain.Store) error {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(store.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ServerError50001
	}
	store.Password = string(cryptedPassword)
	store.Status = domain.StoreStatus.OPEN
	err = su.storeRepository.Create(ctx, store)
	return err
}

func (su *storeUsecase) Signin(ctx context.Context, store *domain.Store) (string, error) {
	var token string
	storeFromDb, err := su.GetByEmail(ctx, store.Email)
	if err != nil {
		return token, domain.ServerError50002
	}
	if storeFromDb == (domain.Store{}) {
		return token, domain.ServerError40402
	}
	err = bcrypt.CompareHashAndPassword([]byte(storeFromDb.Password), []byte(store.Password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return token, domain.ServerError40003
	case err != nil:
		return token, domain.ServerError50001
	}
	token, salt, err := su.generateToken(&storeFromDb)
	if err != nil {
		return token, domain.ServerError50001
	}
	fmt.Println(salt)
	return token, nil
}

type TokenClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func (su *storeUsecase) generateToken(store *domain.Store) (encryptToken string, salt string, err error) {
	now := time.Now()
	claims := TokenClaims{
		store.ID,
		store.Email,
		store.Name,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			// ExpiresAt: now.Add(24 * 30 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	randomUUID := uuid.New().String()
	saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	if err != nil {
		return encryptToken, salt, domain.ServerError50001
	}
	
	encryptToken, err = token.SignedString(saltBytes)
	return encryptToken, string(saltBytes), err
}
