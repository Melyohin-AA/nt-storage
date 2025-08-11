package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func (s *Server) oauthLoginH(w http.ResponseWriter, r *http.Request) {
	oauthState := genStateOauthCookie(w)
	url := s.wad.Config.OauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func genStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(time.Minute * 20)
	buff := make([]byte, 16)
	rand.Read(buff)
	state := base64.URLEncoding.EncodeToString(buff)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func (s *Server) oauthCallbackH(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state")
		return
	}
	token, err := s.wad.Config.OauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err = s.initService(token); err != nil {
		log.Println(err.Error())
		return
	}
	if err = setTokenCookie(w, token); err != nil {
		log.Println(err.Error())
	}
}

func setTokenCookie(w http.ResponseWriter, token *oauth2.Token) error {
	jsonToken, err := json.Marshal(token)
	if err != nil {
		return err
	}
	encoded := base64.URLEncoding.EncodeToString(jsonToken)
	cookie := http.Cookie{Name: "token", Value: encoded, Path: "/", Expires: token.Expiry}
	http.SetCookie(w, &cookie)
	return nil
}

func getTokenCookie(r *http.Request) (*oauth2.Token, error) {
	encoded, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
		return nil, nil
	}
	jsonToken, err := base64.URLEncoding.DecodeString(encoded.Value)
	if err != nil {
		return nil, err
	}
	token := new(oauth2.Token)
	err = json.Unmarshal(jsonToken, token)
	return token, err
}

func (s *Server) redirectToLoginIfRequired(w http.ResponseWriter, r *http.Request) bool {
	if s.wad.Service != nil {
		return false
	}
	token, err := getTokenCookie(r)
	if (err == nil) && (token != nil) {
		err = s.initService(token)
		if err == nil {
			return false
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
	return true
}
