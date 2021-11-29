package httpAPI

import (
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
)

func (had *httpAPIDelivery) customersCreate(w http.ResponseWriter, r *http.Request) {
	// session_id
	// customers - {name, phone, queue_id}
	had.usecase.CreateCustomer(r.Context())
	presenter.JsonResponseOK(w, map[string]string{"hello": "world"})
}
