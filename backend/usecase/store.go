package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
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
	domain            string
}

func NewStoreUsecase(storeRepository domain.StoreRepositoryInterface, signKeyRepository domain.SignKeyRepositoryInterface, logger logging.LoggerTool, domain string) domain.StoreUsecaseInterface {
	return &storeUsecase{storeRepository, signKeyRepository, logger, domain}
}

func (su *storeUsecase) GetByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, err := su.storeRepository.GetByEmail(ctx, email)
	return store, err
}

func (su *storeUsecase) VerifyPasswordLength(password string) error {
	decodedPassword, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	rawPassword := string(decodedPassword)
	// length of password must between 8 and 15.
	if len(rawPassword) < 8 || len(rawPassword) > 15 {
		return domain.ServerError40002
	}
	return nil
}

func (su *storeUsecase) EncryptPassword(password string) (string, error) {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	return string(cryptedPassword), nil
}

func (su *storeUsecase) Create(ctx context.Context, store domain.Store) error {
	store.Status = domain.StoreStatus.OPEN
	err := su.storeRepository.Create(ctx, store)
	return err
}

func (su *storeUsecase) ValidatePassword(ctx context.Context, passwordInDb string, incomingPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInDb), []byte(incomingPassword))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		su.logger.ERRORf("%v", err)
		return domain.ServerError40003
	case err != nil:
		su.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	return nil
}

func (su *storeUsecase) GenerateToken(ctx context.Context, store domain.Store, signKeyType string, expiresDuration time.Duration) (encryptToken string, err error) {
	randomUUID := uuid.New().String()
	saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	signKey := &domain.SignKey{StoreId: store.ID, SignKey: string(saltBytes), SignKeyType: signKeyType}
	err = su.signKeyRepository.Create(ctx, signKey)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := domain.TokenClaims{
		store.ID,
		store.Email,
		store.Name,
		signKey.ID,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expiresDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encryptToken, err = token.SignedString([]byte(signKey.SignKey))
	if err != nil {
		su.logger.ERRORf("%v", err)
		return encryptToken, domain.ServerError50001
	}
	return encryptToken, err
}

func (su *storeUsecase) VerifyToken(ctx context.Context, encryptToken string) (tokenClaims domain.TokenClaims, err error) {
	_, _, err = new(jwt.Parser).ParseUnverified(encryptToken, &tokenClaims)
	if err != nil {
		su.logger.ERRORf("%v", err)
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
	return tokenClaims, nil
}

func (su *storeUsecase) VerifyTokenRenewable(tokenClaims domain.TokenClaims) bool {
	if tokenClaims.ExpiresAt <= time.Now().Add(24*time.Hour).Unix() {
		return true
	}
	return false
}

func (su *storeUsecase) RemoveSignKeyByID(ctx context.Context, signKeyID int) error {
	err := su.signKeyRepository.RemoveByID(ctx, signKeyID)
	return err
}

func (su *storeUsecase) GenerateEmailContentOfForgetPassword(emailToken string, store domain.Store) (subject string, content string) {
	// TODO: update email content to html format.
	resetPasswordUrl := fmt.Sprintf("%s/api/stores/password/renew?emailToken=%s", su.domain, emailToken)
	return "Queue-System Reset Password", fmt.Sprintf("Hello, %s, please click %s", store.Name, resetPasswordUrl)
}

func (su *storeUsecase) Update(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
	err := su.storeRepository.Update(ctx, store, fieldName, newFieldValue)
	return err
}
