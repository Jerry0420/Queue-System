package grpcServices

import (
	context "context"
	"fmt"
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	_ "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
		grpc.WithDefaultServiceConfig(
			`{
				"healthCheckConfig": {
					"serviceName": ""
				}
			}`,
		),
	)
	if err != nil {
		logger.FATALf("grpc connection fail %v", err)
	}
	client := NewGrpcServiceClient(conn)

	healthClient := healthpb.NewHealthClient(conn)
	go func() {
		for {
			resp, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{Service: ""})
			fmt.Println(resp.GetStatus(), "====", err)
			if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {

			}
			time.Sleep(3 * time.Second)
		}
	}()

	return conn, client
}
