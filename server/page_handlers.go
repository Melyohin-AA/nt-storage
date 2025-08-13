package server

import (
	"fmt"
	"net/http"
	"os"
)

type pageHandler struct {
	server  *Server
	content []byte
}

func newPageHandler(server *Server, content string) pageHandler {
	return pageHandler{server: server, content: []byte(content)}
}

func (h pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.server.redirectToLoginIfRequired(w, r) {
		return
	}
	w.Write(h.content)
}

func addPageHandlers(server *Server, mux *http.ServeMux) error {
	var jacket string
	if data, err := os.ReadFile("pages/jacket.html"); err == nil {
		jacket = string(data)
	} else {
		return err
	}
	var index string
	if data, err := os.ReadFile("pages/index.html"); err == nil {
		index = fmt.Sprintf(jacket, string(data))
	} else {
		return err
	}
	mux.Handle("/", newPageHandler(server, index))
	var blobMeta string
	if data, err := os.ReadFile("pages/blob_meta.html"); err == nil {
		blobMeta = fmt.Sprintf(jacket, string(data))
	} else {
		return err
	}
	mux.Handle("/blob/meta", newPageHandler(server, blobMeta))
	var tags string
	if data, err := os.ReadFile("pages/tags.html"); err == nil {
		tags = fmt.Sprintf(jacket, string(data))
	} else {
		return err
	}
	mux.Handle("/tags", newPageHandler(server, tags))
	return nil
}
