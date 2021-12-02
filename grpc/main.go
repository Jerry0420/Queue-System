package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	grpcServices "github.com/jerry0420/queue-system/grpc/proto"
	"github.com/jerry0420/queue-system/grpc/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

type GrpcServicesServer struct {
	grpcServices.UnimplementedGrpcServiceServer
}

func (*GrpcServicesServer) GenerateCSV(ctx context.Context, req *grpcServices.GenerateCSVRequest) (*grpcServices.GenerateCSVResponse, error) {
	fmt.Printf("GenerateCSV function is invoked with %v \n", req)

	name := req.GetName()
	fmt.Println(name)
	content := req.GetContent()
	fmt.Println(content)
	var cotentMap map[string]interface{}
	json.Unmarshal(content, &cotentMap)
	fmt.Println(cotentMap)

	res := &grpcServices.GenerateCSVResponse{
		Filepath: "xxxxxxx",
	}

	return res, nil
}

func (*GrpcServicesServer) SendEmail(ctx context.Context, req *grpcServices.SendEmailRequest) (*grpcServices.SendEmailResponse, error) {
	fmt.Printf("SendEmail function is invoked with %v \n", req)

	subject := req.GetSubject()
	fmt.Println(subject)
	content := req.GetContent()
	fmt.Println(content)
	email := req.GetEmail()
	fmt.Println(email)
	filePath := req.GetFilepath()
	fmt.Println(filePath)

	res := &grpcServices.SendEmailResponse{
		Result: true,
	}

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8002")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}

	opts := []grpc.ServerOption{}

	if config.ServerConfig.ENV() == config.EnvStatus.PROD {
		creds, err := credentials.NewServerTLSFromFile(config.ServerConfig.SERVER_CRT(), config.ServerConfig.SERVER_KEY())
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
		opts = append(opts, grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
				MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
				MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
				Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
			},
		))
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
				PermitWithoutStream: true,            // Allow pings even when there are no active streams
			},
		))
	}

	grpcServer := grpc.NewServer(opts...)
	grpcServices.RegisterGrpcServiceServer(grpcServer, &GrpcServicesServer{})

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}
