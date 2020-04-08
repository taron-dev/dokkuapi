package server

import (
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// Server represents webserver
type Server struct {
	store     *str.Store
	router    http.Handler
	blackList []string
}

// NewServer creates new server
func NewServer(store *str.Store) *Server {
	newServer := new(Server)
	newServer.blackList = []string{}
	newServer.store = store
	newServer.initRoutes()
	newServer.initLogFile()
	return newServer
}

// Start api webserver
func (s *Server) Start() {
	http.ListenAndServe(":3000", s.router)
}
