package httpAPI

import (
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type queueDelivery struct {
	queueUsecase domain.QueueUsecaseInterface
	logger logging.LoggerTool
}

func NewQueueDelivery(router *mux.Router, logger logging.LoggerTool, queueUsecase domain.QueueUsecaseInterface) {
	_ = &queueDelivery{queueUsecase, logger}
}