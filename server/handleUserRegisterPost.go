package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postRegister(w http.ResponseWriter, r *http.Request) {
	handlers.UserRegister(w, r, s.store)
}
