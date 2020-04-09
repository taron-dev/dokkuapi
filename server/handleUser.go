package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.UserRegister(w, r, s.store)
	}
}

func (s *Server) postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.UserLogin(w, r, s.store)
	}
}

func (s *Server) postLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.UserLogout(w, r, &s.blackList)
	}
}

func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.UserDelete(w, r, s.store)
	}
}
