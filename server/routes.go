package server

import (
	"github.com/gorilla/mux"
)

func (s *Server) initRoutes() {
	router := mux.NewRouter()
	router.Handle("/info", s.isAuthenticated(s.getInfo())).Methods("GET")

	router.Handle("/register", s.verifyUser(s.postRegister)).Methods("POST")
	router.Handle("/login", s.verifyUser(s.postLogin)).Methods("POST")
	router.Handle("/logout", s.isAuthenticated(s.postLogout)).Methods("POST")
	router.Handle("/users/{userId}", s.isAuthenticated(s.deleteUser)).Methods("DELETE")

	router.Handle("/apps", s.isAuthenticated(s.postApps)).Methods("POST")
	router.Handle("/apps/{appId}", s.isAuthenticated(s.deleteApp)).Methods("DELETE")
	s.router = router
}
