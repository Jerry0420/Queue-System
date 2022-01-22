package main

import (
	"context"
	"log"
	"net"

	grpcServices "github.com/jerry0420/queue-system/backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcServicesServer struct {
	grpcServices.UnimplementedGrpcServiceServer
}

func (grpcServicesServer *GrpcServicesServer) GenerateCSV(ctx context.Context, req *grpcServices.GenerateCSVRequest) (*grpcServices.GenerateCSVResponse, error) {
	res := &grpcServices.GenerateCSVResponse{
		Filepath: "/im/sample/csv/file/path",
	}

	return res, nil
}

func (grpcServicesServer *GrpcServicesServer) SendEmail(ctx context.Context, req *grpcServices.SendEmailRequest) (*grpcServices.SendEmailResponse, error) {
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

	grpcServer := grpc.NewServer(opts...)
	grpcServicesServer := GrpcServicesServer{}

	grpcServices.RegisterGrpcServiceServer(grpcServer, &grpcServicesServer)

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}