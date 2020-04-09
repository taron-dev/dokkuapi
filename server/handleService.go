package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postAppServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.ServiceCreate(w, r, s.store)
	}
}
