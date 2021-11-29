package httpAPI

import (
	"fmt"
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *httpAPIDelivery) sessionCreate(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := validator.SessionCreate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	tokenClaims, err := had.usecase.VerifyToken(
		r.Context(),
		sessionToken,
		domain.SignKeyTypes.SESSION,
		had.usecase.GetSignKeyByID,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	store := domain.Store{
		ID:    tokenClaims.StoreID,
		Email: tokenClaims.Email,
		Name:  tokenClaims.Name,
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

func (had *httpAPIDelivery) sessionScanned(w http.ResponseWriter, r *http.Request) {
	storeId, sessionId, err := validator.SessionScanned(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	session := domain.StoreSession{
		ID:      sessionId,
		StoreId: int(storeId),
	}
	err = had.usecase.UpdateSession(
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
