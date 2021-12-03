package httpAPI

import (
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) customersCreate(w http.ResponseWriter, r *http.Request) {
	session, customers, err := validator.Customercreate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.CreateCustomer(
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
	// TODO: publish to broker
	presenter.JsonResponseOK(w, customers)
}
