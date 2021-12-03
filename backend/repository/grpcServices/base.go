package grpcServices

import (
	"time"

	"github.com/jerry0420/queue-system/backend/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	// _ "google.golang.org/grpc/health"
	// healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func GetDevGrpcConn(logger logging.LoggerTool, host string) (*grpc.ClientConn, GrpcServiceClient) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		logger.FATALf("grpc connection fail %v", err)
	}
	client := NewGrpcServiceClient(conn)
	
	// healthClient := healthpb.NewHealthClient(conn)
	// go func() {
	// 	for {
	// 		resp, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{Service: ""})
	// 		fmt.Println(resp.GetStatus(), "====", err)
	// 		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {

	// 		}
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	return conn, client
}

func GetGrpcConn(logger logging.LoggerTool, host string, caCrtPath string) (*grpc.ClientConn, GrpcServiceClient) {
	creds, err := credentials.NewClientTLSFromFile(caCrtPath, "queue-system")
	if err != nil {
		logger.FATALf("error: %v", err)
	}
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		},
	))
	opts = append(opts, grpc.WithDefaultServiceConfig(
		`{
			"loadBalancingPolicy":"round_robin",
			"healthCheckConfig": {
				"serviceName": ""
			}
		}`,
	))
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		logger.FATALf("grpc connection fail %v", err)
	}
	client := NewGrpcServiceClient(conn)
	return conn, client
}
