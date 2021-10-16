package usecase

import (
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/domain"
)

type customerUsecase struct {
	customerRepository domain.CustomerRepositoryInterface
	logger logging.LoggerTool
}

func NewCustomerUsecase(customerRepository domain.CustomerRepositoryInterface, logger logging.LoggerTool) domain.CustomerUsecaseInterface {
	return &customerUsecase{customerRepository, logger}
}