package main

import (
	"context"
	"fmt"
	"log"
	"net"

	grpcServices "github.com/jerry0420/queue-system/grpc/proto"
	"google.golang.org/grpc"
)

type GrpcServicesServer struct {
	grpcServices.UnimplementedGrpcServiceServer
}

func (*GrpcServicesServer) GenerateCSV(ctx context.Context, req *grpcServices.GenerateCSVRequest) (*grpcServices.GenerateCSVResponse, error) {
	fmt.Printf("Sum function is invoked with %v \n", req)

	_ = req.GetName()
	_ = req.GetContent()

	res := &grpcServices.GenerateCSVResponse{
		Filepath: "xxxxxxx",
	}

	return res, nil
}

func (*GrpcServicesServer) SendEmail(ctx context.Context, req *grpcServices.SendEmailRequest) (*grpcServices.SendEmailResponse, error) {
	fmt.Printf("Sum function is invoked with %v \n", req)

	_ = req.GetSubject()
	_ = req.GetContent()
	_ = req.GetEmail()
	_ = req.GetFilepath()

	res := &grpcServices.SendEmailResponse{
		Result: true,
	}

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}

	grpcServer := grpc.NewServer()
	grpcServices.RegisterGrpcServiceServer(grpcServer, &GrpcServicesServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}
