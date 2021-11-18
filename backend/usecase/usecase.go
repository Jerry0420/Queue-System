package usecase

import (
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository"
)

type UsecaseConfig struct {
	Domain        string
	StoreDuration time.Duration
}

type usecase struct {
	pgDBRepository repository.RepositoryInterface
	logger         logging.LoggerTool
	config         UsecaseConfig
}

func NewUsecase(
	pgDBRepository repository.RepositoryInterface,
	logger logging.LoggerTool,
	usecaseConfig UsecaseConfig,
) UseCaseInterface {
	return &usecase{pgDBRepository, logger, usecaseConfig}
}