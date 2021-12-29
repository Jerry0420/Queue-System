package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/jerry0420/queue-system/grpc/config"
	grpcServices "github.com/jerry0420/queue-system/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	gomail "gopkg.in/gomail.v2"
)

type GrpcServicesServer struct {
	csvDirPath string
	dialer     *gomail.Dialer
	fromEmail  string
	grpcServices.UnimplementedGrpcServiceServer
}

func (grpcServicesServer *GrpcServicesServer) GenerateCSV(ctx context.Context, req *grpcServices.GenerateCSVRequest) (*grpcServices.GenerateCSVResponse, error) {
	name := req.GetName()
	csvFilePath := filepath.Join(grpcServicesServer.csvDirPath, name+".csv")
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	err = os.Chmod(csvFilePath, 0777)
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(csvFile)

	content := req.GetContent()
	var cotentMap [][]string
	json.Unmarshal(content, &cotentMap)
	err = csvWriter.WriteAll(cotentMap)
	if err != nil {
		return nil, err
	}
	csvWriter.Flush()

	err = csvWriter.Error()
	if err != nil {
		return nil, err
	}

	res := &grpcServices.GenerateCSVResponse{
		Filepath: csvFilePath,
	}

	return res, nil
}

func (grpcServicesServer *GrpcServicesServer) SendEmail(ctx context.Context, req *grpcServices.SendEmailRequest) (*grpcServices.SendEmailResponse, error) {
	message := gomail.NewMessage()
	message.SetHeader("From", grpcServicesServer.fromEmail)

	email := req.GetEmail()
	message.SetHeader("To", email)

	subject := req.GetSubject()
	message.SetHeader("Subject", subject)

	content := req.GetContent()
	message.SetBody("text/html", content)

	filePath := req.GetFilepath()
	if filePath != "" {
		message.Attach(filePath)
	}

	err := grpcServicesServer.dialer.DialAndSend(message)

	if err != nil {
		return nil, err
	}

	res := &grpcServices.SendEmailResponse{
		Result: true,
	}

	if filePath != "" {
		os.Remove(filePath)
	}

	return res, nil
}

func initCsvDirPath(csvDirPath string) error {
	if _, err := os.Stat(csvDirPath); os.IsNotExist(err) {
		err := os.Mkdir(csvDirPath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
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
				MaxConnectionIdle:     15 * time.Second,  // If a client is idle for 15 seconds, send a GOAWAY
				MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
				MaxConnectionAgeGrace: 5 * time.Second,   // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
				Time:                  5 * time.Second,   // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout:               1 * time.Second,   // Wait 1 second for the ping ack before assuming the connection is dead
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
	grpcServicesServer := GrpcServicesServer{
		csvDirPath: "/app/grpc/csvs",
		dialer: gomail.NewDialer(
			config.ServerConfig.EMAIL_SERVER(),
			config.ServerConfig.EMAIL_PORT(),
			config.ServerConfig.EMAIL_USERNAME(),
			config.ServerConfig.EMAIL_PASSWORD(),
		),
		fromEmail: config.ServerConfig.EMAIL_FROM(),
	}
	grpcServices.RegisterGrpcServiceServer(grpcServer, &grpcServicesServer)

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	err = initCsvDirPath(grpcServicesServer.csvDirPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}
