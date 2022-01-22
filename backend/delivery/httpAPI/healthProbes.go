package httpAPI

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/config"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HttpAPIHealthProbes struct {
	db          *sql.DB
	grpcConn    *grpc.ClientConn
	vaultServer string
	env         string
}

func NewHttpAPIHealthProbes(
	router *mux.Router,
	db *sql.DB,
	grpcConn *grpc.ClientConn,
	vaultServer string,
	env string,
) {
	hahp := &HttpAPIHealthProbes{db, grpcConn, vaultServer, env}

	router.HandleFunc(
		"/api/routine/liveness",
		hahp.liveness,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/api/routine/readiness",
		hahp.readiness,
	).Methods(http.MethodGet)

}

func (hahp *HttpAPIHealthProbes) liveness(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (hahp *HttpAPIHealthProbes) readiness(w http.ResponseWriter, r *http.Request) {
	ok := true
	errMsg := ""

	// check database
	err := hahp.db.Ping()
	if err != nil {
		ok = false
		errMsg += err.Error()
		http.Error(w, errMsg, http.StatusServiceUnavailable)
		return
	}

	// check grpc
	healthClient := healthpb.NewHealthClient(hahp.grpcConn)
	resp, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{Service: ""})
	if err != nil {
		ok = false
		errMsg += err.Error()
		http.Error(w, errMsg, http.StatusServiceUnavailable)
		return
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		ok = false
		errMsg += "grpc not ok."
		http.Error(w, errMsg, http.StatusServiceUnavailable)
		return
	}

	// check vault only in production
	if hahp.env == config.EnvStatus.PROD {
		httpClient := http.Client{Timeout: 3 * time.Second}
		response, err := httpClient.Get(hahp.vaultServer + "/v1/sys/health") // url in vault server.
		if err != nil || response.StatusCode != http.StatusOK {
			ok = false
			errMsg += "vault not ok."
			http.Error(w, errMsg, http.StatusServiceUnavailable)
			return
		}
	
		var decodedResponse map[string]interface{}
		json.NewDecoder(response.Body).Decode(&decodedResponse)
		initialized := decodedResponse["initialized"].(bool)
		sealed := decodedResponse["sealed"].(bool)
		standby := decodedResponse["standby"].(bool)
		if initialized == false || sealed == true || standby == true {
			ok = false
			errMsg += "vault not ok."
			http.Error(w, errMsg, http.StatusServiceUnavailable)
			return
		}	
	}

	if ok {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, errMsg, http.StatusServiceUnavailable)
	}
}