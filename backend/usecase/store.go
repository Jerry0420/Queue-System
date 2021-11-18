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
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	err := uc.pgDBRepository.CreateStore(ctx, store, queues)
	return err
}

func (uc *usecase) GetStoreByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, err := uc.pgDBRepository.GetStoreByEmail(ctx, email)
	switch {
	case store != domain.Store{} && err == nil && (time.Now().Sub(store.CreatedAt) < uc.config.StoreDuration):
		return store, domain.ServerError40901
	case store != domain.Store{} && err == nil && (time.Now().Sub(store.CreatedAt) >= uc.config.StoreDuration):
		return store, domain.ServerError40903
	}
	return domain.Store{}, err
}

func (uc *usecase) VerifyPasswordLength(password string) error {
	decodedPassword, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	rawPassword := string(decodedPassword)
	// length of password must between 8 and 15.
	if len(rawPassword) < 8 || len(rawPassword) > 15 {
		return domain.ServerError40002
	}
	return nil
}

func (uc *usecase) EncryptPassword(password string) (string, error) {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	return string(cryptedPassword), nil
}

func (uc *usecase) ValidatePassword(ctx context.Context, passwordInDb string, incomingPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInDb), []byte(incomingPassword))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		uc.logger.ERRORf("%v", err)
		return domain.ServerError40003
	case err != nil:
		uc.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	return nil
}

func (uc *usecase) CloseStore(ctx context.Context, store domain.Store) error {
	// TODO: send report to the store.

	err := uc.pgDBRepository.RemoveStoreByID(ctx, store.ID)
	return err
}

func (uc *usecase) GenerateToken(ctx context.Context, store domain.Store, signKeyType string, expireTime time.Time) (encryptToken string, err error) {
	randomUUID := uuid.New().String()
	saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	signKey := &domain.SignKey{StoreId: store.ID, SignKey: string(saltBytes), SignKeyType: signKeyType}
	err = uc.pgDBRepository.CreateSignKey(ctx, signKey)
	if err != nil {
		return "", err
	}

	claims := domain.TokenClaims{
		store.ID,
		store.Email,
		store.Name,
		store.CreatedAt.Unix(),
		signKey.ID,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encryptToken, err = token.SignedString([]byte(signKey.SignKey))
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return encryptToken, domain.ServerError50001
	}
	return encryptToken, err
}

func (uc *usecase) VerifyToken(ctx context.Context, encryptToken string, signKeyType string, getSignKey func(context.Context, int, string) (domain.SignKey, error)) (tokenClaims domain.TokenClaims, err error) {
	_, _, err = new(jwt.Parser).ParseUnverified(encryptToken, &tokenClaims)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return domain.TokenClaims{}, domain.ServerError40101
	}

	tokenClaims = domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(encryptToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		signKey, err := getSignKey(ctx, tokenClaims.SignKeyID, signKeyType)
		if err != nil {
			return nil, err
		}
		return []byte(signKey.SignKey), nil
	})
	if err != nil {
		uc.logger.ERRORf("%v", err)
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return tokenClaims, domain.ServerError40104
		}
		if serverError, ok := err.(*jwt.ValidationError).Inner.(*domain.ServerError); ok {
			return tokenClaims, serverError
		}
		return tokenClaims, domain.ServerError40103
	}

	if !token.Valid {
		uc.logger.ERRORf("unvalid token")
		return tokenClaims, domain.ServerError40103
	}

	// store expired!
	if time.Now().Sub(time.Unix(tokenClaims.StoreCreatedAt, 0)) >= uc.config.StoreDuration {
		return tokenClaims, domain.ServerError40105
	}

	return tokenClaims, nil
}

func (uc *usecase) RemoveSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error) {
	signKey, err := uc.pgDBRepository.RemoveSignKeyByID(ctx, signKeyID, signKeyType)
	return signKey, err
}

func (uc *usecase) GetSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error) {
	signKey, err := uc.pgDBRepository.GetSignKeyByID(ctx, signKeyID, signKeyType)
	return signKey, err
}

func (uc *usecase) GenerateEmailContentOfForgetPassword(passwordToken string, store domain.Store) (subject string, content string) {
	// TODO: update email content to html format.
	resetPasswordUrl := fmt.Sprintf("%s/stores/password/update?id=%d&password_token=%s", uc.config.Domain, store.ID, passwordToken)
	return "Queue-System Reset Password", fmt.Sprintf("Hello, %s, please click %s", store.Name, resetPasswordUrl)
}

func (uc *usecase) UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
	err := uc.pgDBRepository.UpdateStore(ctx, store, fieldName, newFieldValue)
	return err
}
