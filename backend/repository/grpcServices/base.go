package grpcServices

import (
	"github.com/jerry0420/queue-system/backend/logging"
	"google.golang.org/grpc"
)

func GetDevGrpcConn(logger logging.LoggerTool, host string) (*grpc.ClientConn, *GrpcServiceClient) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		logger.FATALf("grpc connection fail %v", err)
	}
	client := NewGrpcServiceClient(conn)
	return conn, &client
}
