package server

import (
	"encoding/json"
	"io"
	"net/http"
	"nt-storage/storage"
)

func (s *Server) getTagDataH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, s.manager.GetTagData)
}

func (s *Server) addTagCatH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, func() (string, int) {
		var tagCat storage.TagCategory
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err.Error(), http.StatusBadRequest
		}
		if err = json.Unmarshal(body, &tagCat); err != nil {
			return err.Error(), http.StatusBadRequest
		}
		return s.manager.AddTagCat(tagCat)
	})
}

func (s *Server) addTagH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, func() (string, int) {
		var tag storage.Tag
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err.Error(), http.StatusBadRequest
		}
		if err = json.Unmarshal(body, &tag); err != nil {
			return err.Error(), http.StatusBadRequest
		}
		return s.manager.AddTag(tag)
	})
}
