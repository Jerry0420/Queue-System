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
	// 1. check session token
	// 2. flush new session
	// 3. sse, and flush new session

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError50003)
		return
	}

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
