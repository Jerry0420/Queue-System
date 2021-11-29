package grpcServices

import (
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func GetDevGrpcConn(logger logging.LoggerTool, host string) (*grpc.ClientConn, GrpcServiceClient) {
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
				Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
				PermitWithoutStream: true,             // send pings even without active streams
			},
		),
	)
	if err != nil {
		logger.FATALf("grpc connection fail %v", err)
	}
	client := NewGrpcServiceClient(conn)
	return conn, client
}
