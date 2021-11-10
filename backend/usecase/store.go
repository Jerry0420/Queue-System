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
	claims := domain.TokenClaims {
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

func (su *storeUsecase) VerifyToken(ctx context.Context, encryptToken string) (tokenClaims domain.TokenClaims, err error) {
	_, _, err = new(jwt.Parser).ParseUnverified(encryptToken, &tokenClaims)
	if err != nil {
		return domain.TokenClaims{}, domain.ServerError40101
	}

	tokenClaims = domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(encryptToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		signKey, err := su.signKeyRepository.GetByID(ctx, tokenClaims.SignKeyID)
		if err != nil {
			return nil, err
		}
		return []byte(signKey.SignKey), nil
	})
	if err != nil {
		su.logger.ERRORf("%v", err)
		return domain.TokenClaims{}, domain.ServerError40403
	}
	if !token.Valid {
		su.logger.ERRORf("unvalid token")
		return domain.TokenClaims{}, domain.ServerError40101
	}
	if tokenClaims.ExpiresAt <= time.Now().Add(24 * time.Hour).Unix() {
		return tokenClaims, domain.ServerError40103
	}
	return tokenClaims, nil
}

func (su *storeUsecase) RemoveSignKeyByID(ctx context.Context, signKeyID int) error {
	err := su.signKeyRepository.RemoveByID(ctx, signKeyID)
	return err
}