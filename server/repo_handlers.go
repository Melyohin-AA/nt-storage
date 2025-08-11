package server

import (
	"net/http"
)

func (s *Server) initRepoH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, r, s.manager.InitRepo)
}

func (s *Server) pushRepoH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, r, s.manager.PushRepo)
}

func (s *Server) pullRepoH(w http.ResponseWriter, r *http.Request) {
	s.exec(w, r, s.manager.PullRepo)
}
