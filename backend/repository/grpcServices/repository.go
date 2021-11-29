package grpcServices

import (
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
)

type grpcServicesRepository struct {
	client         GrpcServiceClient
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewGrpcServicesRepository(client GrpcServiceClient, logger logging.LoggerTool, contextTimeOut time.Duration) GrpcServicesRepositoryInterface {
	return &grpcServicesRepository{client, logger, contextTimeOut}
}
