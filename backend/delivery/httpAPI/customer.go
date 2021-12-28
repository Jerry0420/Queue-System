package httpAPI

import (
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) customersCreate(w http.ResponseWriter, r *http.Request) {
	session, customers, err := validator.CustomerCreate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.CreateCustomers(
		r.Context(),
		session,
		domain.StoreSessionStatus.SCANNED,
		domain.StoreSessionStatus.USED,
		customers,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	
	go had.broker.Publish(
		had.usecase.TopicNameOfUpdateCustomer(session.StoreId),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, customers)
}

// for used in store...
func (had *httpAPIDelivery) customerUpdate(w http.ResponseWriter, r *http.Request) {
	storeId, oldCustomerStatus, newCustomerStatus, customer, err := validator.CustomerUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.UpdateCustomer(r.Context(), oldCustomerStatus, newCustomerStatus, &customer)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.usecase.TopicNameOfUpdateCustomer(storeId),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, customer)
}