package server

import (
	"net/http"
)

func (s *Server) listFilesH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, r, s.manager.ListFiles)
}

func (s *Server) deleteFileH(w http.ResponseWriter, r *http.Request) {
	fid := r.FormValue("fid")
	s.exec(w, r, func() (string, int) {
		return s.manager.DeleteFile(fid)
	})
}
