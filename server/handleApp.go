package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) postApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppsCreate(w, r, s.store)
	}
}

func (s *Server) getApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppsGet(w, r, s.store)
	}
}

func (s *Server) deleteApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppDelete(w, r, s.store)
	}
}

func (s *Server) putApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppEdit(w, r, s.store)
	}
}

func (s *Server) postAppDeploy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppDeploy(w, r)
	}
}

func (s *Server) putAppStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AppStop(w, r)
	}
}
