package httpAPI

import (
	"errors"
	"net/http"
	"time"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) storeOpen(w http.ResponseWriter, r *http.Request) {
	store, queues, err := validator.StoreOpen(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.VerifyPasswordLength(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	encryptedPassword, err := had.usecase.EncryptPassword(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store.Password = encryptedPassword

	storeInDb, err := had.usecase.GetStoreByEmail(r.Context(), store.Email)
	switch {
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		err = had.usecase.CloseStore(r.Context(), storeInDb)
		if err != nil {
			presenter.JsonResponse(w, nil, err)
			return
		}
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40901):
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.CreateStore(r.Context(), &store, queues)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreWithQueuesForResponse(store, queues))
}

func (had *httpAPIDelivery) storeSignin(w http.ResponseWriter, r *http.Request) {
	incomingStore, err := validator.StoreSignin(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	storeInDb, err := had.usecase.GetStoreByEmail(r.Context(), incomingStore.Email)
	switch {
	case storeInDb == domain.Store{} && err != nil:
		presenter.JsonResponse(w, nil, err)
		return
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		_ = had.usecase.CloseStore(r.Context(), storeInDb)
		presenter.JsonResponse(w, nil, domain.ServerError40903)
		return
	}

	err = had.usecase.ValidatePassword(r.Context(), storeInDb.Password, incomingStore.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	refreshTokenExpiresAt := storeInDb.CreatedAt.Add(had.config.StoreDuration)
	token, err := had.usecase.GenerateToken(
		r.Context(),
		storeInDb,
		domain.SignKeyTypes.REFRESH,
		refreshTokenExpiresAt,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	cookie := http.Cookie{
		Name:     domain.SignKeyTypes.REFRESH,
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     V_1("/stores/token"),
		MaxAge:   int(refreshTokenExpiresAt.Sub(time.Now())),
	}
	http.SetCookie(w, &cookie)
	presenter.JsonResponseOK(w, presenter.StoreForResponse(storeInDb))
}

func (had *httpAPIDelivery) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	encryptedRefreshToken, err := validator.StoreTokenRefresh(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	tokenClaims, err := had.usecase.VerifyToken(
		r.Context(),
		encryptedRefreshToken.Value,
		domain.SignKeyTypes.REFRESH,
		had.usecase.GetSignKeyByID,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	tokenExpiresAt := time.Now().Add(had.config.TokenDuration)
	// normal token
	normalToken, err := had.usecase.GenerateToken(
		r.Context(),
		store,
		domain.SignKeyTypes.NORMAL,
		tokenExpiresAt,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	// session token
	sessionToken, err := had.usecase.GenerateToken(
		r.Context(),
		store,
		domain.SignKeyTypes.SESSION,
		tokenExpiresAt,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreToken(store, normalToken, tokenExpiresAt, sessionToken))
}

func (had *httpAPIDelivery) storeClose(w http.ResponseWriter, r *http.Request) {
	tokenClaims, err := validator.StoreClose(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	err = had.usecase.CloseStore(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *httpAPIDelivery) passwordForgot(w http.ResponseWriter, r *http.Request) {
	store, err := validator.StorePasswordForgot(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err = had.usecase.GetStoreByEmail(r.Context(), store.Email)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	passwordToken, err := had.usecase.GenerateToken(r.Context(), store, domain.SignKeyTypes.PASSWORD, time.Now().Add(had.config.PasswordTokenDuration))
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	_, _ = had.usecase.GenerateEmailContentOfForgetPassword(passwordToken, store)
	// TODO: SendEmail function (grpc)
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *httpAPIDelivery) passwordUpdate(w http.ResponseWriter, r *http.Request) {
	jsonBody, id, err := validator.StorePasswordUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	tokenClaims, err := had.usecase.VerifyToken(r.Context(), jsonBody["password_token"], domain.SignKeyTypes.PASSWORD, had.usecase.RemoveSignKeyByID)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	if id != tokenClaims.StoreID {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}

	err = had.usecase.VerifyPasswordLength(jsonBody["password"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		Password:  jsonBody["password"],
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	encryptedPassword, err := had.usecase.EncryptPassword(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store.Password = encryptedPassword

	err = had.usecase.UpdateStore(r.Context(), &store, "password", store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}
