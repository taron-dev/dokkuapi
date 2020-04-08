package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postApps(w http.ResponseWriter, r *http.Request) {
	handlers.AppsCreate(w, r, s.store)
}

func (s *Server) deleteApp(w http.ResponseWriter, r *http.Request) {
	handlers.AppDelete(w, r, s.store)
}
