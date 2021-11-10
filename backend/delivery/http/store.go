package delivery

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/middleware"
	"github.com/jerry0420/queue-system/backend/presenter"
)

type storeDelivery struct {
	storeUsecase domain.StoreUsecaseInterface
	logger       logging.LoggerTool
}

func NewStoreDelivery(router *mux.Router, mw *middleware.Middleware, logger logging.LoggerTool, storeUsecase domain.StoreUsecaseInterface) {
	sd := &storeDelivery{storeUsecase, logger}
	router.HandleFunc(
		"/stores/signup",
		sd.signup,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		"/stores/signin",
		sd.signin,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.Handle(
		"/stores/signout",
		mw.AuthenticationMiddleware(http.HandlerFunc(sd.signout)),
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")
}

func (sd *storeDelivery) signup(w http.ResponseWriter, r *http.Request) {
	var store domain.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil || store.Name == "" || store.Email == "" || store.Password == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(store.Password)
	rawPassword := string(decodedPassword)
	// length of password must between 8 and 15.
	if err != nil || len(rawPassword) < 8 || len(rawPassword) > 15 {
		presenter.JsonResponse(w, nil, domain.ServerError40002)
		return
	}
	err = sd.storeUsecase.Create(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}

func (sd *storeDelivery) signin(w http.ResponseWriter, r *http.Request) {
	var store domain.Store
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil || store.Email == "" || store.Password == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}

	store, err = sd.storeUsecase.Signin(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	token, err := sd.storeUsecase.GenerateToken(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, map[string]interface{}{"token": token})
}

func (sd *storeDelivery) signout(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(domain.TokenClaims)

	var jsonBody map[string]int
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil || jsonBody["storeID"] != token.StoreID || jsonBody["signKeyID"] != token.SignKeyID {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}

	err = sd.storeUsecase.RemoveSignKeyByID(r.Context(), token.SignKeyID)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}
