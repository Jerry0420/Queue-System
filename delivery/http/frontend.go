package delivery

import (
	"io/fs"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/logging"
)

type frontendDelivery struct {
	logger logging.LoggerTool
	frontendFiles fs.FS
}

func NewFrontendDelivery(router *mux.Router, logger logging.LoggerTool, frontendFiles fs.FS) {
	fd := &frontendDelivery{logger: logger, frontendFiles: frontendFiles}
	router.PathPrefix("/").HandlerFunc(fd.serveFrontend)
}

func (fd *frontendDelivery) serveFrontend(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(fd.frontendFiles)).ServeHTTP(w, r)
}