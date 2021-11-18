package httpAPI

import (
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/usecase"
)

type customerDelivery struct {
	usecase usecase.UseCaseInterface
	logger  logging.LoggerTool
}

func NewCustomerDelivery(router *mux.Router, logger logging.LoggerTool, usecase usecase.UseCaseInterface) {
	_ = &customerDelivery{usecase, logger}
}
