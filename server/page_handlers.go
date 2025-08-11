package server

import (
	"fmt"
	"net/http"
	"os"
)

type pageHandler struct {
	content []byte
}

func newPageHandler(content string) pageHandler {
	return pageHandler{content: []byte(content)}
}

func (h pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.content)
}

func addPageHandlers(mux *http.ServeMux) error {
	var template string
	if data, err := os.ReadFile("pages/template.html"); err == nil {
		template = string(data)
	} else {
		return err
	}
	var index string
	if data, err := os.ReadFile("pages/index.html"); err == nil {
		index = fmt.Sprintf(template, string(data))
	} else {
		return err
	}
	mux.Handle("/", newPageHandler(index))
	var blobMeta string
	if data, err := os.ReadFile("pages/blob_meta.html"); err == nil {
		blobMeta = fmt.Sprintf(template, string(data))
	} else {
		return err
	}
	mux.Handle("/blob/meta", newPageHandler(blobMeta))
	return nil
}
