package httpAPI

import (
	"fmt"
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) createSession(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := validator.SessionCreate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err := had.usecase.VerifySessionToken(r.Context(), sessionToken)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError50003)
		return
	}
	consumerChan := had.broker.Subscribe(had.usecase.TopicNameOfUpdateSession(store.ID))
	defer had.broker.UnsubscribeConsumer(had.usecase.TopicNameOfUpdateSession(store.ID), consumerChan)

	session, err := had.usecase.CreateSession(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	fmt.Fprintf(w, "data: %v\n\n", presenter.SessionCreate(had.config.Domain, session))
	flusher.Flush()

	for {
		select {
		case event := <-consumerChan:
			if event["old_session_id"].(string) == session.ID {
				session, err = had.usecase.CreateSession(r.Context(), store)
				if err != nil {
					presenter.JsonResponse(w, nil, err)
					return
				}
				fmt.Fprintf(w, "data: %v\n\n", presenter.SessionCreate(had.config.Domain, session))
				flusher.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (had *httpAPIDelivery) scannedSession(w http.ResponseWriter, r *http.Request) {
	session, err := validator.SessionScanned(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.usecase.UpdateSessionStatus(
		r.Context(),
		&session,
		domain.StoreSessionStatus.NORMAL,  //oldStatus
		domain.StoreSessionStatus.SCANNED, //newStatus
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.usecase.TopicNameOfUpdateSession(session.StoreId),
		map[string]interface{}{"old_session_id": session.ID},
	)
	presenter.JsonResponseOK(w, session)
}
