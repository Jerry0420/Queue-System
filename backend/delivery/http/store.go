package delivery

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/presenter"
)

type storeDelivery struct {
	storeUsecase domain.StoreUsecaseInterface
	logger       logging.LoggerTool
}

func NewStoreDelivery(router *mux.Router, logger logging.LoggerTool, storeUsecase domain.StoreUsecaseInterface) {
	sd := &storeDelivery{storeUsecase, logger}
	router.HandleFunc("/stores", sd.create).Methods(http.MethodPost)
}

func (sd *storeDelivery) create(w http.ResponseWriter, r *http.Request) {
	presenter.JsonResponseOK(w, map[string]string{"hello": "world"})
}
