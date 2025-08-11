package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nt-storage/manager"
	"nt-storage/storage"
	"sync"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Server struct {
	wad     *storage.Wad
	manager manager.Manager
	server  *http.Server
	mx      sync.RWMutex
}

func NewServer(wad *storage.Wad) (*Server, error) {
	s := &Server{
		wad:     wad,
		manager: manager.NewManager(wad),
		server:  &http.Server{Addr: fmt.Sprintf(":%d", wad.Config.Port)},
	}
	wad.Config.OauthConfig.RedirectURL = fmt.Sprintf("http://localhost:%d/auth/cb", wad.Config.Port)
	var err error
	s.server.Handler, err = s.newHandlers()
	return s, err
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); (err != nil) && (err != http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) newHandlers() (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	if err := addPageHandlers(mux); err != nil {
		return mux, err
	}
	mux.HandleFunc("/auth/login", s.oauthLoginH)
	mux.HandleFunc("/auth/cb", s.oauthCallbackH)
	mux.HandleFunc("/api/file/list", s.listFilesH)
	mux.HandleFunc("/api/file/delete", s.deleteFileH)
	mux.HandleFunc("/api/repo/init", s.initRepoH)
	mux.HandleFunc("/api/repo/push", s.pushRepoH)
	mux.HandleFunc("/api/repo/pull", s.pullRepoH)
	mux.HandleFunc("/api/blob/list", s.listBlobsH)
	mux.HandleFunc("/api/blob/add", s.addBlobH)
	mux.HandleFunc("/api/blob/update", s.updateBlobH)
	mux.HandleFunc("/api/blob/fetch", s.fetchBlobH)
	mux.HandleFunc("/api/blob/delete", s.deleteBlobH)
	mux.HandleFunc("/api/blob/meta", s.blobMetaH)
	return mux, nil
}

func (s *Server) initService(token *oauth2.Token) error {
	httpClient := s.wad.Config.OauthConfig.Client(context.Background(), token)
	service, err := drive.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}
	s.wad.Service = service
	return nil
}

func (s *Server) exec(w http.ResponseWriter, r *http.Request, f func() (string, int)) {
	if s.redirectToLoginIfRequired(w, r) {
		return
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	output, code := f()
	if code/100 != 2 {
		if code/100 == 5 {
			log.Println(output)
		}
		http.Error(w, output, code)
		return
	}
	fmt.Fprint(w, output)
}
