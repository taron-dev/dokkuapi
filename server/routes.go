package server

import (
	"github.com/gorilla/mux"
)

func (s *Server) initRoutes() {
	router := mux.NewRouter()
	router.Handle("/info", s.isAuthenticated(s.getInfo())).Methods("GET")

	router.Handle("/register", s.verifyUser(s.postRegister())).Methods("POST")
	router.Handle("/login", s.verifyUser(s.postLogin())).Methods("POST")
	router.HandleFunc("/logout", s.isAuthenticated(s.postLogout())).Methods("POST")
	router.Handle("/users/{userId}", s.isAuthenticated(s.deleteUser())).Methods("DELETE")

	router.Handle("/apps", s.isAuthenticated(s.postApps())).Methods("POST")
	router.Handle("/apps", s.isAuthenticated(s.getApps())).Methods("GET")
	router.Handle("/apps/{appId}", s.isAuthenticated(s.isUserAuthorizedApp(s.deleteApp()))).Methods("DELETE")
	router.Handle("/apps/{appId}", s.isAuthenticated(s.isUserAuthorizedApp(s.putApp()))).Methods("PUT")

	router.Handle("/apps/{appId}/deploy", s.isAuthenticated(s.isUserAuthorizedApp(s.postAppDeploy()))).Methods("POST")
	router.Handle("/apps/{appId}/stop", s.isAuthenticated(s.isUserAuthorizedApp(s.putAppStop()))).Methods("PUT")

	router.Handle("/apps/{appId}/instances", s.isAuthenticated(s.isUserAuthorizedApp(s.getAppInstances()))).Methods("GET")
	router.Handle("/apps/{appId}/instances", s.isAuthenticated(s.isUserAuthorizedApp(s.putAppInstances()))).Methods("PUT")

	router.Handle("/apps/{appId}/services", s.isAuthenticated(s.isUserAuthorizedApp(s.postAppServices()))).Methods("POST")
	router.Handle("/apps/{appId}/services", s.isAuthenticated(s.isUserAuthorizedApp(s.getAppServices()))).Methods("GET")
	router.Handle("/apps/{appId}/services/{serviceId}", s.isAuthenticated(s.isUserAuthorizedApp(s.deleteAppService()))).Methods("DELETE")

	s.router = router
}
