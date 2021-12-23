package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jerry0420/queue-system/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

func (uc *Usecase) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
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

func (uc *Usecase) SigninStore(ctx context.Context, email string, password string) (store domain.Store, token string, refreshTokenExpiresAt time.Time, err error) {
	store, err = uc.pgDBRepository.GetStoreByEmail(ctx, email)
	store, err = uc.CheckStoreExpirationStatus(store, err)
	switch {
	case store == domain.Store{} && err != nil:
		return store, token, refreshTokenExpiresAt, err
	case store != domain.Store{} && errors.Is(err, domain.ServerError40903):
		err = uc.CloseStore(ctx, store)
		return store, token, refreshTokenExpiresAt, err
	}

	err = uc.ValidatePassword(store.Password, password)
	if err != nil {
		return store, token, refreshTokenExpiresAt, err
	}

	refreshTokenExpiresAt = store.CreatedAt.Add(uc.config.StoreDuration)
	token, err = uc.GenerateToken(
		ctx,
		store,
		domain.SignKeyTypes.REFRESH,
		refreshTokenExpiresAt,
	)
	if err != nil {
		return store, token, refreshTokenExpiresAt, err
	}

	return store, token, refreshTokenExpiresAt, nil
}

func (uc *Usecase) RefreshToken(ctx context.Context, encryptedRefreshToken string) (
	store domain.Store,
	normalToken string,
	sessionToken string,
	tokenExpiresAt time.Time,
	err error,
) {
	tokenClaims, err := uc.VerifyToken(
		ctx,
		encryptedRefreshToken,
		domain.SignKeyTypes.REFRESH,
		true,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}
	store = domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}

	tokenExpiresAt = time.Now().Add(uc.config.TokenDuration)
	// normal token
	normalToken, err = uc.GenerateToken(
		ctx,
		store,
		domain.SignKeyTypes.NORMAL,
		tokenExpiresAt,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}
	// session token
	sessionToken, err = uc.GenerateToken(
		ctx,
		store,
		domain.SignKeyTypes.SESSION,
		tokenExpiresAt,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}

	return store, normalToken, sessionToken, tokenExpiresAt, nil
}

func (uc *Usecase) CloseStore(ctx context.Context, store domain.Store) error {
	// TODO: send report to the store.

	err := uc.pgDBRepository.RemoveStoreByID(ctx, store.ID)
	return err
}

func (uc *Usecase) ForgetPassword(ctx context.Context, email string) (store domain.Store, err error) {
	store, err = uc.pgDBRepository.GetStoreByEmail(ctx, email)
	if err != nil {
		return store, err
	}
	passwordToken, err := uc.GenerateToken(
		ctx,
		store,
		domain.SignKeyTypes.PASSWORD,
		time.Now().Add(uc.config.PasswordTokenDuration),
	)
	if err != nil {
		return store, err
	}

	_, _ = uc.GenerateEmailContentOfForgetPassword(passwordToken, store)
	// TODO: SendEmail function (grpc)

	return store, err
}

func (uc *Usecase) UpdatePassword(ctx context.Context, passwordToken string, newPassword string) (store domain.Store, err error) {
	tokenClaims, err := uc.VerifyToken(
		ctx,
		passwordToken,
		domain.SignKeyTypes.PASSWORD,
		false,
	)
	if err != nil {
		return store, err
	}
	store = domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}

	encryptedPassword, err := uc.EncryptPassword(newPassword)
	if err != nil {
		return store, err
	}

	err = uc.pgDBRepository.UpdateStore(ctx, &store, "password", encryptedPassword)
	if err != nil {
		return store, err
	}

	store.Password = encryptedPassword
	return store, nil
}

func (uc *Usecase) UpdateStoreDescription(ctx context.Context, newDescription string, store *domain.Store) error {
	err := uc.pgDBRepository.UpdateStore(ctx, store, "description", newDescription)
	if err != nil {
		return err
	}
	return nil
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

func (uc *Usecase) VerifyToken(ctx context.Context, encryptToken string, signKeyType string, withSignkeyPreserved bool) (tokenClaims domain.TokenClaims, err error) {
	_, _, err = new(jwt.Parser).ParseUnverified(encryptToken, &tokenClaims)
	if err != nil {
		uc.logger.ERRORf("%v", err)
		return domain.TokenClaims{}, domain.ServerError40101
	}

	tokenClaims = domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(encryptToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		var getSignKeyFunc func(context.Context, int, string) (domain.SignKey, error)
		if withSignkeyPreserved == true {
			getSignKeyFunc = uc.pgDBRepository.GetSignKeyByID
		} else {
			getSignKeyFunc = uc.pgDBRepository.RemoveSignKeyByID
		}
		signKey, err := getSignKeyFunc(ctx, tokenClaims.SignKeyID, signKeyType)
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

func (uc *Usecase) GenerateEmailContentOfForgetPassword(passwordToken string, store domain.Store) (subject string, content string) {
	// TODO: update email content to html format.
	resetPasswordUrl := fmt.Sprintf("%s/stores/%d/password/update?password_token=%s", uc.config.Domain, store.ID, passwordToken)
	return "Queue-System Reset Password", fmt.Sprintf("Hello, %s, please click %s", store.Name, resetPasswordUrl)
}

func (uc *Usecase) TopicNameOfUpdateCustomer(storeId int) string {
	return fmt.Sprintf("updateCustomer.%d", storeId)
}

func (uc *Usecase) GetStoreWIthQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error) {
	store, err := uc.pgDBRepository.GetStoreWIthQueuesAndCustomersById(ctx, storeId)
	return store, err
}

func (uc *Usecase) VerifyNormalToken(ctx context.Context, normalToken string) (tokenClaims domain.TokenClaims, err error) {
	encryptToken := strings.Split(normalToken, " ")
	if len(encryptToken) == 2 && strings.ToLower(encryptToken[0]) == "bearer" {
		tokenClaims, err = uc.VerifyToken(
			ctx,
			encryptToken[1],
			domain.SignKeyTypes.NORMAL,
			true,
		)
		return tokenClaims, err
	}
	return tokenClaims, domain.ServerError40102
}

func (uc *Usecase) VerifySessionToken(ctx context.Context, sessionToken string) (store domain.Store, err error) {
	tokenClaims, err := uc.VerifyToken(
		ctx,
		sessionToken,
		domain.SignKeyTypes.SESSION,
		true, // TODO: change to RemoveSignKeyByID
	)
	if err != nil {
		return store, err
	}
	store = domain.Store{
		ID:    tokenClaims.StoreID,
		Email: tokenClaims.Email,
		Name:  tokenClaims.Name,
	}
	return store, nil
}
