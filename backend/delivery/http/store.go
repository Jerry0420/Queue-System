package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/middleware"
	"github.com/jerry0420/queue-system/backend/presenter"
)

type StoreDeliveryConfig struct {
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
}

type storeDelivery struct {
	logger       logging.LoggerTool
	storeUsecase domain.StoreUsecaseInterface
	config       StoreDeliveryConfig
}

func NewStoreDelivery(router *mux.Router, logger logging.LoggerTool, mw *middleware.Middleware, storeUsecase domain.StoreUsecaseInterface, storeDeliveryConfig StoreDeliveryConfig) {
	sd := &storeDelivery{logger, storeUsecase, storeDeliveryConfig}
	router.HandleFunc(
		V_1("/stores"),
		sd.open,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/signin"),
		sd.signin,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/token"),
		sd.tokenRefresh,
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(sd.close)),
	).Methods(http.MethodDelete)

	router.HandleFunc(
		V_1("/stores/password/forgot"),
		sd.passwordForgot,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/password"),
		sd.passwordUpdate,
	).Methods(http.MethodPatch).Headers("Content-Type", "application/json")
}

func (sd *storeDelivery) open(w http.ResponseWriter, r *http.Request) {
	var store domain.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil || store.Name == "" || store.Email == "" || store.Password == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	err = sd.storeUsecase.VerifyPasswordLength(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	encryptedPassword, err := sd.storeUsecase.EncryptPassword(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store.Password = encryptedPassword

	storeInDb, err := sd.storeUsecase.GetByEmail(r.Context(), store.Email)
	switch {
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		err = sd.storeUsecase.Close(r.Context(), storeInDb)
		if err != nil {
			presenter.JsonResponse(w, nil, err)
			return
		}
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40901):
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = sd.storeUsecase.Create(r.Context(), &store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (sd *storeDelivery) signin(w http.ResponseWriter, r *http.Request) {
	var incomingStore domain.Store
	err := json.NewDecoder(r.Body).Decode(&incomingStore)
	if err != nil || incomingStore.Email == "" || incomingStore.Password == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}

	storeInDb, err := sd.storeUsecase.GetByEmail(r.Context(), incomingStore.Email)
	switch {
	case storeInDb == domain.Store{} && err != nil:
		presenter.JsonResponse(w, nil, err)
		return
	case storeInDb != domain.Store{} && errors.Is(err, domain.ServerError40903):
		_ = sd.storeUsecase.Close(r.Context(), storeInDb)
		presenter.JsonResponse(w, nil, domain.ServerError40903)
		return
	}

	err = sd.storeUsecase.ValidatePassword(r.Context(), storeInDb.Password, incomingStore.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	refreshTokenExpiresAt := storeInDb.CreatedAt.Add(sd.config.StoreDuration)
	token, err := sd.storeUsecase.GenerateToken(
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

func (sd *storeDelivery) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	encryptedRefreshToken, err := r.Cookie(domain.SignKeyTypes.REFRESH)
	if err != nil || len(encryptedRefreshToken.Value) == 0 {
		presenter.JsonResponse(w, nil, domain.ServerError40102)
		return
	}
	tokenClaims, err := sd.storeUsecase.VerifyToken(
		r.Context(),
		encryptedRefreshToken.Value,
		domain.SignKeyTypes.REFRESH,
		sd.storeUsecase.GetSignKeyByID,
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
	tokenExpiresAt := time.Now().Add(sd.config.TokenDuration)
	token, err := sd.storeUsecase.GenerateToken(
		r.Context(),
		store,
		domain.SignKeyTypes.NORMAL,
		tokenExpiresAt,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreToken(store, token, tokenExpiresAt))
}

func (sd *storeDelivery) close(w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value(domain.SignKeyTypes.NORMAL).(domain.TokenClaims)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id != tokenClaims.StoreID {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	err = sd.storeUsecase.Close(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (sd *storeDelivery) passwordForgot(w http.ResponseWriter, r *http.Request) {
	var store domain.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil || store.Email == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}

	store, err = sd.storeUsecase.GetByEmail(r.Context(), store.Email)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	passwordToken, err := sd.storeUsecase.GenerateToken(r.Context(), store, domain.SignKeyTypes.PASSWORD, time.Now().Add(sd.config.PasswordTokenDuration))
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	_, _ = sd.storeUsecase.GenerateEmailContentOfForgetPassword(passwordToken, store)
	// TODO: SendEmail function (grpc)
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (sd *storeDelivery) passwordUpdate(w http.ResponseWriter, r *http.Request) {
	var jsonBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil || jsonBody["passwordToken"] == "" || jsonBody["password"] == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	tokenClaims, err := sd.storeUsecase.VerifyToken(r.Context(), jsonBody["passwordToken"], domain.SignKeyTypes.PASSWORD, sd.storeUsecase.RemoveSignKeyByID)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	err = sd.storeUsecase.VerifyPasswordLength(jsonBody["password"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id != tokenClaims.StoreID {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		Password:  jsonBody["password"],
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	encryptedPassword, err := sd.storeUsecase.EncryptPassword(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store.Password = encryptedPassword

	err = sd.storeUsecase.Update(r.Context(), &store, "password", store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}
