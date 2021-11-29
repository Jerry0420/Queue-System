package grpcServices

import context "context"

type GrpcServicesRepositoryInterface interface {
	// services.go
	GenerateCSV(ctx context.Context, name string, content []byte) (filePath string, err error)
	SendEmail(ctx context.Context, subject string, content string, email string, filepath string) (result bool, err error)
}
