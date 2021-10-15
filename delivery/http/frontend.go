package delivery

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/logging"
)

type frontendDelivery struct {
	logger logging.LoggerTool
	contentStatic fs.FS
	baseDir string
}

func NewFrontendDelivery(router *mux.Router, logger logging.LoggerTool, baseDir string, contentStatic fs.FS) {
	fd := &frontendDelivery{logger: logger, contentStatic: contentStatic, baseDir: baseDir}
	router.PathPrefix("/").HandlerFunc(fd.serveFrontend)
}

func (fd *frontendDelivery) serveFrontend(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(fd.baseDir, path)
	_, err = os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.FS(fd.contentStatic)).ServeHTTP(w, r)
}