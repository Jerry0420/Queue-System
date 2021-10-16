package delivery

import (
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type storeDelivery struct {
	storeUsecase domain.StoreUsecaseInterface
	logger logging.LoggerTool
}

func NewStoreDelivery(router *mux.Router, logger logging.LoggerTool, storeUsecase domain.StoreUsecaseInterface) {
	_ = &storeDelivery{storeUsecase, logger}
}