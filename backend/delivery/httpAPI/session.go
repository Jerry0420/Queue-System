package httpAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/domain"
)

// var (
// 	counter int
// )

func (had *httpAPIDelivery) SessionCreate(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.URL.Query().Get("session_token") // Because, it's a GET method.
	if sessionToken == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
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

	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError50003)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}

	session, err := had.usecase.CreateSession(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	var flushedData bytes.Buffer
	// TODO: combine url and qrcode
	json.NewEncoder(&flushedData).Encode(session)
	fmt.Fprintf(w, "data: %v\n\n", flushedData.String())
	flusher.Flush()
	flushedData.Reset()

	consumerChan := had.broker.Subscribe("updateSession." + strconv.Itoa(store.ID))
	defer had.broker.UnsubscribeConsumer("updateSession."+strconv.Itoa(store.ID), consumerChan)

	for {
		event := <-consumerChan
		if event["old_session_id"].(string) == session.ID {
			session, err = had.usecase.CreateSession(r.Context(), store)
			if err != nil {
				presenter.JsonResponse(w, nil, err)
				return
			}
			// TODO: combine url and qrcode
			json.NewEncoder(&flushedData).Encode(session)
			fmt.Fprintf(w, "data: %v\n\n", flushedData.String())
			flusher.Flush()
			flushedData.Reset()
		}
		// select {
		// case event := <-consumerChan:
		// 	if event["old_session_id"].(string) == session.ID {
		// 		fmt.Println(event)
		// 		session, err := had.usecase.CreateSession(r.Context(), store)
		// 		if err != nil {
		// 			presenter.JsonResponse(w, nil, err)
		// 			return
		// 		}
		// 		// TODO: combine url and qrcode
		// 		json.NewEncoder(&flushedData).Encode(session)
		// 		fmt.Fprintf(w, "data: %v\n\n", flushedData.String())
		// 		flusher.Flush()
		// 		flushedData.Reset()
		// 	}
		// case <-r.Context().Done():
		// 	fmt.Println("close~~~~~~")
		// 	return
		// }
	}

	// t := time.NewTicker(2 * time.Second)
	// defer t.Stop()
	// for {
	// 	select {
	// 	case <-t.C:
	// 		counter++
	// 		c := counter
	// 		var buf bytes.Buffer
	// 		json.NewEncoder(&buf).Encode(map[string]interface{}{"hello": c})
	// 		fmt.Fprintf(w, "data: %v\n\n", buf.String())
	// 		flusher.Flush()
	// 	case <-r.Context().Done():
	// 		return
	// 	}
	// }

	// var buf bytes.Buffer
	// json.NewEncoder(&buf).Encode(map[string]interface{}{"hello": "world"})
	// fmt.Fprintf(w, "data: %v\n\n", buf.String())
	// flusher.Flush()

	// buf.Reset()
	// json.NewEncoder(&buf).Encode(map[string]interface{}{"hello": "2222"})
	// fmt.Fprintf(w, "data: %v\n\n", buf.String())
	// flusher.Flush()
}

func (had *httpAPIDelivery) SessionScan(w http.ResponseWriter, r *http.Request) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	storeId, ok := jsonBody["store_id"].(float64)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	vars := mux.Vars(r)
	sessionId, ok := vars["id"]
	if !ok || sessionId == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40004)
		return
	}
	session := domain.StoreSession{
		ID:      sessionId,
		StoreId: int(storeId),
	}
	oldStatus := domain.StoreSessionStatus.NORMAL
	newStatus := domain.StoreSessionStatus.SCANNED
	err = had.usecase.UpdateSession(r.Context(), &session, oldStatus, newStatus)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	go had.broker.Publish("updateSession."+strconv.Itoa(session.StoreId), map[string]interface{}{"old_session_id": session.ID})
	presenter.JsonResponseOK(w, session)
}
