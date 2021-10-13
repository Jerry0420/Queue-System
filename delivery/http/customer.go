package delivery

import (
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/domain"
)

type customerDelivery struct {
	customerUsecase domain.CustomerUsecaseInterface
	logger logging.LoggerTool
}

func NewCustomerDelivery(router *mux.Router, logger logging.LoggerTool, customerUsecase domain.CustomerUsecaseInterface) {
	_ = &customerDelivery{customerUsecase, logger}
}