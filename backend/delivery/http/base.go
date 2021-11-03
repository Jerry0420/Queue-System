package delivery

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/presenter"
)

type baseDelivery struct {
	logger logging.LoggerTool
}

func NewBaseDelivery(router *mux.Router, logger logging.LoggerTool) {
	bd := &baseDelivery{logger}
	// these two routes will just response to the client directly, and will not go into any middleware. 
	router.MethodNotAllowedHandler = http.HandlerFunc(bd.methodNotAllow)
	router.NotFoundHandler = http.HandlerFunc(bd.notFound)
}

func (bd *baseDelivery) methodNotAllow(w http.ResponseWriter, r *http.Request) {
	presenter.JsonResponse(w, nil, domain.ServerError40501)
}

func (bd *baseDelivery) notFound(w http.ResponseWriter, r *http.Request) {
	presenter.JsonResponse(w, nil, domain.ServerError40401)
}