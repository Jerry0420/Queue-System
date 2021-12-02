package validator

import (
	"encoding/json"
	"net/http"

	"github.com/jerry0420/queue-system/backend/domain"
)

func Customercreate(r *http.Request) (storeId int, sessionId string, customers []domain.Customer, err error) {
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return storeId, sessionId, customers, domain.ServerError40001
	}
	
	sessionId, ok := jsonBody["session_id"].(string)
	if !ok || sessionId == "" {
		return storeId, sessionId, customers, domain.ServerError40001
	}
	
	storeIdFloat64, ok := jsonBody["store_id"].(float64)
	if !ok {
		return storeId, sessionId, customers, domain.ServerError40001
	}
	
	customersInfo, ok := jsonBody["customers"].([]interface{})
	if !ok || len(customersInfo) <= 0 {
		return storeId, sessionId, customers, domain.ServerError40001
	}
	
	for _, customerInfo := range customersInfo {
		customerInfo, ok := customerInfo.(map[string]interface{})
		if !ok {
			return storeId, "", []domain.Customer{}, domain.ServerError40001
		}
		name, ok := customerInfo["name"].(string)
		if !ok || name == "" {
			return storeId, "", []domain.Customer{}, domain.ServerError40001
		}
		phone, ok := customerInfo["phone"].(string)
		if !ok || phone == "" {
			return storeId, "", []domain.Customer{}, domain.ServerError40001
		}
		queueId, ok := customerInfo["queue_id"].(float64)
		if !ok {
			return storeId, "", []domain.Customer{}, domain.ServerError40001
		}
		customers = append(customers, domain.Customer{Name: name, Phone: phone, QueueID: int(queueId), Status: domain.CustomerStatus.NORMAL})
	}

	return int(storeIdFloat64), sessionId, customers, nil
}
