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

func (uc *Usecase) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	err := uc.VerifyPasswordLength(store.Password)
	if err != nil {
		return err
	}
	encryptedPassword, err := uc.EncryptPassword(store.Password)
	if err != nil {
		return err
	}
	store.Password = encryptedPassword

	storeInDb, err := uc.pgDBRepository.GetStoreByEmail(ctx, store.Email)
	storeInDb, err = uc.CheckStoreExpirationStatus(storeInDb, err)
	switch {
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		err = uc.CloseStore(ctx, storeInDb)
		if err != nil {
			return err
		}
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40901):
		return err
	}

	err = uc.pgDBRepository.CreateStore(ctx, store, queues)
	return err
}

func (uc *Usecase) SigninStore(ctx context.Context, incomingStore *domain.Store) (token string, refreshTokenExpiresAt time.Time,err error) {
	storeInDb, err := uc.pgDBRepository.GetStoreByEmail(ctx, incomingStore.Email)
	switch {
	case storeInDb == domain.Store{} && err != nil:
		return token, refreshTokenExpiresAt, err
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		err = uc.CloseStore(ctx, storeInDb)
		return token, refreshTokenExpiresAt, err
	}

	err = uc.ValidatePassword(storeInDb.Password, incomingStore.Password)
	if err != nil {
		return token, refreshTokenExpiresAt, err
	}

	refreshTokenExpiresAt = storeInDb.CreatedAt.Add(uc.config.StoreDuration)
	token, err = uc.GenerateToken(
		ctx,
		storeInDb,
		domain.SignKeyTypes.REFRESH,
		refreshTokenExpiresAt,
	)
	if err != nil {
		return token, refreshTokenExpiresAt, err
	}
	
	*incomingStore = storeInDb
	return token, refreshTokenExpiresAt, nil
}

func (uc *Usecase) RefreshToken(ctx context.Context, store *domain.Store) (
	normalToken string, 
	sessionToken string, 
	tokenExpiresAt time.Time, 
	err error,
	) {
	tokenExpiresAt = time.Now().Add(uc.config.TokenDuration)
	// normal token
	normalToken, err = uc.GenerateToken(
		ctx,
		*store,
		domain.SignKeyTypes.NORMAL,
		tokenExpiresAt,
	)
	if err != nil {
		return normalToken, sessionToken, tokenExpiresAt, err
	}
	// session token
	sessionToken, err = uc.GenerateToken(
		ctx,
		*store,
		domain.SignKeyTypes.SESSION,
		tokenExpiresAt,
	)
	if err != nil {
		return normalToken, sessionToken, tokenExpiresAt, err
	}

	return normalToken, sessionToken, tokenExpiresAt, nil
}

func (uc *Usecase) GetStoreByEmail(ctx context.Context, email string) (domain.Store, error) {
	store, err := uc.pgDBRepository.GetStoreByEmail(ctx, email)
	store, err = uc.CheckStoreExpirationStatus(store, err)
	return store, err
}

func (uc *Usecase) CheckStoreExpirationStatus(store domain.Store, err error) (domain.Store, error) {
	switch {
	case store != domain.Store{} && err == nil && (time.Now().Sub(store.CreatedAt) < uc.config.StoreDuration):
		return store, domain.ServerError40901
	case store != domain.Store{} && err == nil && (time.Now().Sub(store.CreatedAt) >= uc.config.StoreDuration):
		return store, domain.ServerError40903
	}
	return domain.Store{}, err
}

func (uc *Usecase) GetStoreWIthQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error) {
	store, err := uc.pgDBRepository.GetStoreWIthQueuesAndCustomersById(ctx, storeId)
	return store, err
}

func (uc *Usecase) VerifyPasswordLength(password string) error {
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

func (uc *Usecase) EncryptPassword(password string) (string, error) {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	return string(cryptedPassword), nil
}

func (uc *Usecase) ValidatePassword(passwordInDb string, incomingPassword string) error {
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

func (uc *Usecase) CloseStore(ctx context.Context, store domain.Store) error {
	// TODO: send report to the store.

	err := uc.pgDBRepository.RemoveStoreByID(ctx, store.ID)
	return err
}

func (uc *Usecase) GenerateToken(ctx context.Context, store domain.Store, signKeyType string, expireTime time.Time) (encryptToken string, err error) {
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

func (uc *Usecase) VerifyToken(ctx context.Context, encryptToken string, signKeyType string, getSignKey func(context.Context, int, string) (domain.SignKey, error)) (tokenClaims domain.TokenClaims, err error) {
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

func (uc *Usecase) RemoveSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error) {
	signKey, err := uc.pgDBRepository.RemoveSignKeyByID(ctx, signKeyID, signKeyType)
	return signKey, err
}

func (uc *Usecase) GetSignKeyByID(ctx context.Context, signKeyID int, signKeyType string) (domain.SignKey, error) {
	signKey, err := uc.pgDBRepository.GetSignKeyByID(ctx, signKeyID, signKeyType)
	return signKey, err
}

func (uc *Usecase) GenerateEmailContentOfForgetPassword(passwordToken string, store domain.Store) (subject string, content string) {
	// TODO: update email content to html format.
	resetPasswordUrl := fmt.Sprintf("%s/stores/%d/password/update?password_token=%s", uc.config.Domain, store.ID, passwordToken)
	return "Queue-System Reset Password", fmt.Sprintf("Hello, %s, please click %s", store.Name, resetPasswordUrl)
}

func (uc *Usecase) UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error {
	err := uc.pgDBRepository.UpdateStore(ctx, store, fieldName, newFieldValue)
	return err
}

func (uc *Usecase) TopicNameOfUpdateCustomer(storeId int) string {
	return fmt.Sprintf("updateCustomer.%d", storeId)
}