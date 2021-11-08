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
		return token, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storeFromDb.Password), []byte(store.Password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return token, domain.ServerError40003
	case err != nil:
		return token, domain.ServerError50001
	}
	signKey, err := su.generateSalt(ctx, storeFromDb.ID)
	if err != nil {
		return token, err
	}
	token, err = su.generateToken(store, signKey)
	return token, err
}

func (su *storeUsecase) generateSalt(ctx context.Context, storeId int) (signKey *domain.SignKey, err error) {
	randomUUID := uuid.New().String()
	saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	if err != nil {
		return &domain.SignKey{}, domain.ServerError50001
	}
	signKey = &domain.SignKey{StoreId: storeId, SignKey: string(saltBytes), SignKeyType: domain.SignKeyTypes.SIGNIN}
	signKeyID, err := su.signKeyRepository.Create(ctx, signKey)
	if err != nil {
		return &domain.SignKey{}, err
	}
	signKey.ID = signKeyID
	return signKey, nil
}

type tokenClaims struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	SignKeyID int    `json:"signkey_id"`
	jwt.StandardClaims
}

func (su *storeUsecase) generateToken(store *domain.Store, signKey *domain.SignKey) (encryptToken string, err error) {
	now := time.Now()
	claims := tokenClaims{
		store.ID,
		store.Email,
		store.Name,
		signKey.ID,
		jwt.StandardClaims{
			IssuedAt: now.Unix(),
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
