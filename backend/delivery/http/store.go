package delivery

import (
	"encoding/json"
	"net/http"
	"time"

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
	).Methods(http.MethodDelete).Headers("Content-Type", "application/json")

	router.HandleFunc(
		"/stores/password/forget",
		sd.passwordForget,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		"/stores/password/update",
		sd.passwordUpdate,
	).Methods(http.MethodPatch).Headers("Content-Type", "application/json")
}

func (sd *storeDelivery) signup(w http.ResponseWriter, r *http.Request) {
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
	sd.logger.ERRORf(encryptedPassword)
	store.Password = encryptedPassword

	err = sd.storeUsecase.Create(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}

func (sd *storeDelivery) signin(w http.ResponseWriter, r *http.Request) {
	var incomingStore domain.Store
	err := json.NewDecoder(r.Body).Decode(&incomingStore)
	if err != nil || incomingStore.Email == "" || incomingStore.Password == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	storeInDb, err := sd.storeUsecase.GetByEmail(r.Context(), incomingStore.Email)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	
	err = sd.storeUsecase.ValidatePassword(r.Context(), storeInDb.Password, incomingStore.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	incomingStore = storeInDb

	token, err := sd.storeUsecase.GenerateToken(r.Context(), incomingStore, domain.SignKeyTypes.SIGNIN, 24*time.Hour)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, map[string]interface{}{"token": token})
}

func (sd *storeDelivery) signout(w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("token").(domain.TokenClaims)
	var jsonBody map[string]int
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil || jsonBody["id"] != tokenClaims.StoreID {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}

	err = sd.storeUsecase.RemoveSignKeyByID(r.Context(), tokenClaims.SignKeyID)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}

func (sd *storeDelivery) passwordForget(w http.ResponseWriter, r *http.Request) {
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
	emailToken, err := sd.storeUsecase.GenerateToken(r.Context(), store, domain.SignKeyTypes.EMAIL, 5*time.Minute)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	// TODO: SendEmail function (grpc)
	_, content := sd.storeUsecase.GenerateEmailContentOfForgetPassword(emailToken, store)
	// TODO: return nil
	presenter.JsonResponseOK(w, map[string]string{"emailContent": content})
}

func (sd *storeDelivery) passwordUpdate(w http.ResponseWriter, r *http.Request) {
	var jsonBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil || jsonBody["emailToken"] == "" || jsonBody["password"] == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	tokenClaims, err := sd.storeUsecase.VerifyToken(r.Context(), jsonBody["emailToken"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	err = sd.storeUsecase.VerifyPasswordLength(jsonBody["password"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store := domain.Store{ID: tokenClaims.StoreID, Email: tokenClaims.Email, Name: tokenClaims.Name, Password: jsonBody["password"]}
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

	err = sd.storeUsecase.RemoveSignKeyByID(r.Context(), tokenClaims.SignKeyID)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}
