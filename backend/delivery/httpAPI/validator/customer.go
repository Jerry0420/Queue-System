package validator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
)

func CustomerCreate(r *http.Request) (session domain.StoreSession, customers []domain.Customer, err error) {
	session = r.Context().Value(domain.StoreSessionString).(domain.StoreSession)

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return session, customers, domain.ServerError40001
	}

	storeIdFloat64, ok := jsonBody["store_id"].(float64)
	if !ok {
		return session, customers, domain.ServerError40001
	}

	customersInfo, ok := jsonBody["customers"].([]interface{})
	if !ok || len(customersInfo) <= 0 {
		return session, customers, domain.ServerError40001
	}
// exceed 5 customers at a req is forbidden.
	if len(customersInfo) > 5 {
		return session, customers, domain.ServerError40005
	}

	for _, customerInfo := range customersInfo {
		customerInfo, ok := customerInfo.(map[string]interface{})
		if !ok {
			return session, customers, domain.ServerError40001
		}
		name, ok := customerInfo["name"].(string)
		if !ok || name == "" {
			return session, customers, domain.ServerError40001
		}
		phone, ok := customerInfo["phone"].(string)
		if !ok || phone == "" {
			return session, customers, domain.ServerError40001
		}
		queueId, ok := customerInfo["queue_id"].(float64)
		if !ok {
			return session, customers, domain.ServerError40001
		}
		customers = append(
			customers, 
			domain.Customer{Name: name, Phone: phone, QueueID: int(queueId), Status: domain.CustomerStatus.NORMAL},
		)
	}

	if int(storeIdFloat64) != session.StoreId {
		return session, customers, domain.ServerError40004
	}

	return session, customers, nil
}

func CustomerUpdate(r *http.Request) (storeId int, oldCustomerStatus string, newCustomerStatus string, customer domain.Customer, err error) {
	tokenClaims := r.Context().Value(domain.SignKeyTypes.NORMAL).(domain.TokenClaims)

	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["id"])
	if err != nil {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	storeIdFloat64, ok := jsonBody["store_id"].(float64)
	if !ok {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	if int(storeIdFloat64) != tokenClaims.StoreID {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40004
	}

	queueIdFloat64, ok := jsonBody["queue_id"].(float64)
	if !ok {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	oldCustomerStatus, ok = jsonBody["old_customer_status"].(string)
	if !ok || oldCustomerStatus == "" {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	newCustomerStatus, ok = jsonBody["new_customer_status"].(string)
	if !ok || newCustomerStatus == "" {
		return storeId, oldCustomerStatus, newCustomerStatus, customer, domain.ServerError40001
	}

	customer = domain.Customer{
		ID: customerId,
		QueueID: int(queueIdFloat64),
		Status: newCustomerStatus,
	}

	return int(storeIdFloat64), oldCustomerStatus, newCustomerStatus, customer, nil
}
