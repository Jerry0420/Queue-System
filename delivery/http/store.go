package delivery

import (
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
)

type storeDelivery struct {
	storeUsecase domain.StoreUsecaseInterface
	logger logging.LoggerTool
}

func NewStoreDelivery(router *mux.Router, logger logging.LoggerTool, storeUsecase domain.StoreUsecaseInterface) {
	_ = &storeDelivery{storeUsecase, logger}
}