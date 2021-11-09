package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"golang.org/x/crypto/bcrypt"
)

type storeUsecase struct {
	storeRepository   domain.StoreRepositoryInterface
	signKeyRepository domain.SignKeyRepositoryInterface
	logger            logging.LoggerTool
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, signKeyRepository domain.SignKeyRepositoryInterface, logger logging.LoggerTool) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, signKeyRepository, logger}
}

func (su *storeUsecase) GetByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, serverError := su.storeRepository.GetByEmail(ctx, email)
	return store, serverError
}

func (su *storeUsecase) Create(ctx context.Context, store domain.Store) error {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(store.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ServerError50001
	}
	store.Password = string(cryptedPassword)
	store.Status = domain.StoreStatus.OPEN
	err = su.storeRepository.Create(ctx, store)
	return err
}

func (su *storeUsecase) Signin(ctx context.Context, store domain.Store) (domain.Store, error) {
	storeFromDb, err := su.GetByEmail(ctx, store.Email)
	if err != nil {
		return domain.Store{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storeFromDb.Password), []byte(store.Password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return domain.Store{}, domain.ServerError40003
	case err != nil:
		return domain.Store{}, domain.ServerError50001
	}
	return storeFromDb, nil
}

type tokenClaims struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	SignKeyID int    `json:"signkey_id"`
	jwt.StandardClaims
}

func (su *storeUsecase) GenerateToken(ctx context.Context, store domain.Store) (encryptToken string, err error) {
	randomUUID := uuid.New().String()
	saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.ServerError50001
	}
	signKey := &domain.SignKey{StoreId: store.ID, SignKey: string(saltBytes), SignKeyType: domain.SignKeyTypes.SIGNIN}
	err = su.signKeyRepository.Create(ctx, signKey)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := tokenClaims{
		store.ID,
		store.Email,
		store.Name,
		signKey.ID,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(24 * 30 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encryptToken, err = token.SignedString([]byte(signKey.SignKey))
	if err != nil {
		return encryptToken, domain.ServerError50001
	}
	return encryptToken, err
}

func (su *storeUsecase) ValidateToken(ctx context.Context, encryptToken string) (store domain.Store, err error) {
	var claims tokenClaims
	_, _, err = new(jwt.Parser).ParseUnverified(encryptToken, &claims)
	if err != nil {
		return domain.Store{}, domain.ServerError40101
	}

	claims = tokenClaims{}
	token, err := jwt.ParseWithClaims(encryptToken, &claims, func(token *jwt.Token) (interface{}, error) {
		signKey, err := su.signKeyRepository.GetByID(ctx, claims.SignKeyID)
		if err != nil {
			return nil, err
		}
		return []byte(signKey.SignKey), nil
	})
	if err != nil {
		return domain.Store{}, domain.ServerError40101
	}
	if !token.Valid {
		return domain.Store{}, domain.ServerError40101
	}
	store = domain.Store{ID: claims.ID, Email: claims.Email, Name: claims.Name}
	return store, nil
}
