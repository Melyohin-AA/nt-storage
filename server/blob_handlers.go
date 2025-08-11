package server

import (
	"encoding/json"
	"io"
	"net/http"
	"nt-storage/manager"
	"nt-storage/storage"
)

func (s *Server) listBlobsH(w http.ResponseWriter, r *http.Request) {
	var filter manager.Filter
	filter.Read(r)
	s.exec(w, r, func() (string, int) {
		return s.manager.ListBlobs(filter)
	})
}

func (s *Server) addBlobH(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	s.exec(w, r, func() (string, int) {
		return s.manager.AddBlob(name)
	})
}

func (s *Server) updateBlobH(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	s.exec(w, r, func() (string, int) {
		return s.manager.UpdateBlob(fid)
	})
}

func (s *Server) fetchBlobH(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	s.exec(w, r, func() (string, int) {
		return s.manager.FetchBlob(fid)
	})
}

func (s *Server) deleteBlobH(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	s.exec(w, r, func() (string, int) {
		return s.manager.DeleteBlob(fid)
	})
}

func (s *Server) blobMetaH(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	switch r.Method {
	case http.MethodGet:
		s.exec(w, r, func() (string, int) {
			return s.manager.GetBlobMeta(fid)
		})
	case http.MethodPut:
		s.exec(w, r, func() (string, int) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err.Error(), http.StatusBadRequest
			}
			var meta storage.Blob
			if err = json.Unmarshal(body, &meta); err != nil {
				return err.Error(), http.StatusBadRequest
			}
			return s.manager.SetBlobMeta(fid, meta)
		})
	default:
		http.Error(w, "suppoted methods: GET, PUT", http.StatusMethodNotAllowed)
	}
}
