package server

import (
	"github.com/gorilla/mux"
	//"github.com/ondro2208/dokkuapi/controller"
)

func (s *Server) initRoutes() {
	router := mux.NewRouter()
	//router.Handle("/info", s.getInfo()).Methods("GET")
	router.Handle("/info", s.IsAuthenticated(s.getInfo())).Methods("GET")

	router.Handle("/register", s.verifyUser(s.postRegister)).Methods("POST")
	router.Handle("/login", s.verifyUser(s.postLogin)).Methods("POST")
	router.Handle("/logout", s.IsAuthenticated(s.postLogout)).Methods("POST")
	router.Handle("/users/{userId}", s.IsAuthenticated(s.deleteUser)).Methods("DELETE")
	router.Handle("/apps", s.IsAuthenticated(s.postApps)).Methods("POST")
	//router.Handle("/apps/{appId}", s.IsAuthenticated(controller.DeleteApp)).Methods("DELETE")
	s.router = router
}
