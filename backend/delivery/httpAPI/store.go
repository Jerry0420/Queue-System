package httpAPI

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) openStore(w http.ResponseWriter, r *http.Request) {
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

	err = had.usecase.CreateStore(r.Context(), &store, queues)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreWithQueuesForResponse(store, queues))
}

func (had *httpAPIDelivery) signinStore(w http.ResponseWriter, r *http.Request) {
	incomingStore, err := validator.StoreSignin(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, token, refreshTokenExpiresAt, err := had.usecase.SigninStore(
		r.Context(),
		incomingStore.Email,
		incomingStore.Password,
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
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *httpAPIDelivery) refreshToken(w http.ResponseWriter, r *http.Request) {
	encryptedRefreshToken, err := validator.StoreTokenRefresh(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, normalToken, sessionToken, tokenExpiresAt, err := had.usecase.RefreshToken(r.Context(), encryptedRefreshToken.Value)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreToken(store, normalToken, tokenExpiresAt, sessionToken))
}

func (had *httpAPIDelivery) closeStore(w http.ResponseWriter, r *http.Request) {
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

func (had *httpAPIDelivery) closeStorerRoutine(w http.ResponseWriter, r *http.Request) {
	deletedStoresCount, err := had.usecase.CloseStoreRoutine(r.Context())
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, map[string]interface{}{"result": deletedStoresCount})
}

func (had *httpAPIDelivery) forgotPassword(w http.ResponseWriter, r *http.Request) {
	store, err := validator.StorePasswordForgot(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err = had.usecase.ForgetPassword(r.Context(), store.Email)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *httpAPIDelivery) updatePassword(w http.ResponseWriter, r *http.Request) {
	jsonBody, _, err := validator.StorePasswordUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.VerifyPasswordLength(jsonBody["password"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err := had.usecase.UpdatePassword(
		r.Context(), 
		jsonBody["password_token"], 
		jsonBody["password"],
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *httpAPIDelivery) getStoreInfo(w http.ResponseWriter, r *http.Request) {
	storeId, err := validator.StoreInfoGet(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError50003)
		return
	}
	consumerChan := had.broker.Subscribe(had.usecase.TopicNameOfUpdateCustomer(storeId))
	defer had.broker.UnsubscribeConsumer(had.usecase.TopicNameOfUpdateCustomer(storeId), consumerChan)

	store, err := had.usecase.GetStoreWIthQueuesAndCustomersById(r.Context(), storeId)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	fmt.Fprintf(w, "data: %v\n\n", presenter.StoreGet(store))
	flusher.Flush()

	for {
		select {
		case <-consumerChan:
			store, err := had.usecase.GetStoreWIthQueuesAndCustomersById(r.Context(), storeId)
			if err != nil {
				presenter.JsonResponse(w, nil, err)
				return
			}
			fmt.Fprintf(w, "data: %v\n\n", presenter.StoreGet(store))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (had *httpAPIDelivery) updateStoreDescription(w http.ResponseWriter, r *http.Request) {
	store, err := validator.StoreDescriptionUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.UpdateStoreDescription(r.Context(), store.Description, &store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.usecase.TopicNameOfUpdateCustomer(store.ID),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}
