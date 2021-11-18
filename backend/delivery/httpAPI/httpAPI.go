package httpAPI

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/usecase"
)

type HttpAPIDeliveryConfig struct {
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
}

type httpAPIDelivery struct {
	logger  logging.LoggerTool
	usecase usecase.UseCaseInterface
	config  HttpAPIDeliveryConfig
}

func NewHttpAPIDelivery(router *mux.Router, logger logging.LoggerTool, mw *middleware.Middleware, usecase usecase.UseCaseInterface, config HttpAPIDeliveryConfig) {
	had := &httpAPIDelivery{logger, usecase, config}

	// stores
	router.HandleFunc(
		V_1("/stores"),
		had.storeOpen,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/signin"),
		had.storeSignin,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/token"),
		had.tokenRefresh,
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.storeClose)),
	).Methods(http.MethodDelete)

	router.HandleFunc(
		V_1("/stores/password/forgot"),
		had.passwordForgot,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/password"),
		had.passwordUpdate,
	).Methods(http.MethodPatch).Headers("Content-Type", "application/json")

	//queues

	//customers

	// base routes
	// these two routes will just response to the client directly, and will not go into any middleware.
	router.MethodNotAllowedHandler = http.HandlerFunc(had.methodNotAllow)
	router.NotFoundHandler = http.HandlerFunc(had.notFound)
}
