package usecase

import (
	"context"
	"encoding/json"
	"fmt"
)

// TODO: remove！！！！
func (uc *usecase) CreateCustomer(ctx context.Context) {
	contentOfCSV, _ := json.Marshal(map[string]string{
		"hello": "world",
	})
	filePath, err := uc.grpcServicesRepository.GenerateCSV(ctx, "jerry_hospital", contentOfCSV)
	fmt.Println(filePath, err)

	result, err := uc.grpcServicesRepository.SendEmail(ctx, "subjectxxx", "contentxxx", "emailxxx", "filepathxxx")
	fmt.Println(result, err)
}
