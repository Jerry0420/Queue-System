package usecase

import (
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type queueUsecase struct {
	queueRepository domain.QueueRepositoryInterface
	logger logging.LoggerTool
}

func NewQueueUsecase(queueRepository domain.QueueRepositoryInterface, logger logging.LoggerTool) domain.QueueUsecaseInterface {
	return &queueUsecase{queueRepository, logger}
}