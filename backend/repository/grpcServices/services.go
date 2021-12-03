package grpcServices

import (
	context "context"

	"github.com/jerry0420/queue-system/backend/domain"
)

func (repo *grpcServicesRepository) GenerateCSV(ctx context.Context, name string, content []byte) (filePath string, err error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	req := &GenerateCSVRequest{
		Name:    name,
		Content: content,
	}
	res, err := repo.client.GenerateCSV(ctx, req)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return filePath, domain.ServerError50004
	}

	return res.GetFilepath(), nil
}

func (repo *grpcServicesRepository) SendEmail(ctx context.Context, subject string, content string, email string, filepath string) (result bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, repo.contextTimeOut)
	defer cancel()

	req := &SendEmailRequest{
		Subject:  subject,
		Content:  content,
		Email:    email,
		Filepath: filepath,
	}
	res, err := repo.client.SendEmail(ctx, req)
	if err != nil {
		repo.logger.ERRORf("error %v", err)
		return result, domain.ServerError50004
	}

	return res.GetResult(), nil
}
