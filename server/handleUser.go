package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postRegister(w http.ResponseWriter, r *http.Request) {
	handlers.UserRegister(w, r, s.store)
}

func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	handlers.UserLogin(w, r, s.store)
}

func (s *Server) postLogout(w http.ResponseWriter, r *http.Request) {
	handlers.UserLogout(w, r, &s.blackList)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	handlers.UserDelete(w, r, s.store)
}
