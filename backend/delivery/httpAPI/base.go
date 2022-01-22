package httpAPI

import (
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *HttpAPIDelivery) methodNotAllow(w http.ResponseWriter, r *http.Request) {
	presenter.JsonResponse(w, nil, domain.ServerError40501)
}

func (had *HttpAPIDelivery) notFound(w http.ResponseWriter, r *http.Request) {
	presenter.JsonResponse(w, nil, domain.ServerError40401)
}
