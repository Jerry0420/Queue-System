package usecase

import (
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/grpcServices"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
)

type UsecaseConfig struct {
	Domain                string
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
	GrpcReplicaCount      int
}

type Usecase struct {
	pgDBRepository         *pgDB.PgDBRepository
	grpcServicesRepository grpcServices.GrpcServicesRepositoryInterface
	logger                 logging.LoggerTool
	config                 UsecaseConfig
}

func NewUsecase(
	pgDBRepository *pgDB.PgDBRepository,
	grpcServicesRepository grpcServices.GrpcServicesRepositoryInterface,
	logger logging.LoggerTool,
	usecaseConfig UsecaseConfig,
) *Usecase {
	return &Usecase{pgDBRepository, grpcServicesRepository, logger, usecaseConfig}
}
