package httpAPI

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/usecase"
)

type HttpAPIDeliveryConfig struct {
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
	Domain                string
}

type httpAPIDelivery struct {
	logger             logging.LoggerTool
	customerUsecase    usecase.CustomerUseCaseInterface
	sessionUsecase     usecase.SessionUseCaseInterface
	storeUsecase       usecase.StoreUseCaseInterface
	integrationUsecase usecase.IntegrationUseCaseInterface
	broker             *broker.Broker
	config             HttpAPIDeliveryConfig
}

func NewHttpAPIDelivery(
	router *mux.Router,
	logger logging.LoggerTool,
	mw *middleware.Middleware,
	customerUsecase usecase.CustomerUseCaseInterface,
	sessionUsecase usecase.SessionUseCaseInterface,
	storeUsecase usecase.StoreUseCaseInterface,
	integrationUsecase usecase.IntegrationUseCaseInterface,
	broker *broker.Broker,
	config HttpAPIDeliveryConfig,
) {
	had := &httpAPIDelivery{logger, customerUsecase, sessionUsecase, storeUsecase, integrationUsecase, broker, config}

	// stores
	router.HandleFunc(
		V_1("/stores"),
		had.openStore,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/signin"),
		had.signinStore,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/token"),
		had.refreshToken,
	).Methods(http.MethodPut)

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.closeStore)),
	).Methods(http.MethodDelete)

	router.HandleFunc(
		V_1("/stores/password/forgot"),
		had.forgotPassword,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/password"),
		had.updatePassword,
	).Methods(http.MethodPatch).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/sse"),
		had.getStoreInfoWithSSE,
	).Methods(http.MethodGet) // get method for sse.

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}"),
		had.getStoreInfo,
	).Methods(http.MethodGet)

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.updateStoreDescription)),
	).Methods(http.MethodPut)

	router.HandleFunc(
		V_1("/routine/stores"),
		had.closeStorerRoutine,
	).Methods(http.MethodDelete)

	//queues

	// sessions
	router.HandleFunc(
		V_1("/sessions/sse"),
		had.createSession,
	).Methods(http.MethodGet) // get method for sse.

	router.Handle(
		V_1("/sessions/{id}"),
		mw.SessionAuthenticationMiddleware(http.HandlerFunc(had.scannedSession)),
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	//customers
	router.Handle(
		V_1("/customers"),
		mw.SessionAuthenticationMiddleware(http.HandlerFunc(had.customersCreate)),
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.Handle(
		V_1("/customers/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.customerUpdate)),
	).Methods(http.MethodPut)

	// base routes
	// these two routes will just response to the client directly, and will not go into any middleware.
	router.MethodNotAllowedHandler = http.HandlerFunc(had.methodNotAllow)
	router.NotFoundHandler = http.HandlerFunc(had.notFound)
}
