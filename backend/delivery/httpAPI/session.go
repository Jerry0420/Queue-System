package httpAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/domain"
)

var (
	counter int
)

func (had *httpAPIDelivery) SessionCreate(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.URL.Query().Get("session_token")
	if sessionToken == "" {
		presenter.JsonResponse(w, nil, domain.ServerError40001)
		return
	}
	tokenClaims, err := had.usecase.VerifyToken(r.Context(), sessionToken, domain.SignKeyTypes.SESSION, had.usecase.GetSignKeyByID)
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
	json.NewEncoder(&flushedData).Encode(session)
	fmt.Fprintf(w, "data: %v\n\n", flushedData.String())
	flusher.Flush()
	flushedData.Reset()

	// TODO: sub, flush new session in for loop
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			counter++
			c := counter
			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(map[string]interface{}{"hello": c})
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}

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
	// update session status
	// notify sessionCreate to flush new session to clients
}
