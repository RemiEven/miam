package rest

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"time"
)

//go:embed static/*
var staticContentFs embed.FS

const (
	index      = "index.html"
	baseFolder = "static/"
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

	content, err := staticContentFs.ReadFile(baseFolder + upath)
	switch {
	case err == nil:
		http.ServeContent(responseWriter, request, upath, applicationStartTime, bytes.NewReader(content))
	case errors.Is(err, fs.ErrNotExist):
		serveIndex(responseWriter, request)
	default:
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
	}
}

func serveIndex(responseWriter http.ResponseWriter, request *http.Request) {
	if content, err := staticContentFs.ReadFile(baseFolder + index); err != nil {
		http.Error(responseWriter, "Internal server error", http.StatusInternalServerError)
	} else {
		http.ServeContent(responseWriter, request, index, applicationStartTime, bytes.NewReader(content))
	}
}
