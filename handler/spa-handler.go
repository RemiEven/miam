package handler

import (
	"net/http"
	"os"
	"time"
)

const (
	index      = "index.html"
	baseFolder = "./static/"
)

var applicationStartTime time.Time

func init() {
	applicationStartTime = time.Now()
}

// SpaHandler is a handler for serving Single Page Application
type SpaHandler struct{}

func (handler SpaHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	upath := request.URL.Path
	// Path sanitizing is already done by mux; no need for further cleaning
	if upath == "/" || upath == "." {
		upath = index
	}

	file, err := os.Open(baseFolder + upath)
	switch {
	case err == nil:
		http.ServeContent(responseWriter, request, upath, applicationStartTime, file)
	case os.IsNotExist(err):
		serveIndex(responseWriter, request)
	default:
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
	}
}

func serveIndex(responseWriter http.ResponseWriter, request *http.Request) {
	if file, err := os.Open(baseFolder + index); err != nil {
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
	} else {
		http.ServeContent(responseWriter, request, index, applicationStartTime, file)
	}
}
