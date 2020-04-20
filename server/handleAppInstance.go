package server

import (
	"github.com/ondro2208/dokkuapi/handlers"
	"net/http"
)

func (s *Server) getAppInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.InstancesGet(w, r)
	}
}
