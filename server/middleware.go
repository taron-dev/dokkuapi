package server

import (
	auth "github.com/ondro2208/dokkuapi/authentication"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"net/http"
	"strings"
)

// IsAuthenticated verifies if request is authenticated
func (s *Server) IsAuthenticated(endpointHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("content-type", "application/json")
			if auth.HasValidToken(response, request, s.blackList) {
				// database usage for user id ??? TODO
				request = contextimpl.DecorateWithSub(request)
				endpointHandler(response, request)
			} else {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte("Not Authorized"))
			}
		})
}

// VerifyUser handles user with access token
func (s *Server) verifyUser(endpointHandler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		accessToken := request.Header.Get("Authorization")
		accessToken = strings.Split(accessToken, "Bearer ")[1]
		githubUser, err := auth.GetGithubUser(accessToken)
		if err != nil {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not Authorized"))
			return
		}
		// database usage for github user id ??? TODO
		request = contextimpl.DecorateWithGithubUser(request, githubUser)
		endpointHandler(response, request)
	})
}
