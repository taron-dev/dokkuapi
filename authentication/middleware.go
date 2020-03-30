package authentication

import (
	"context"
	"net/http"
	"strings"
)

// IsAuthenticated verifies if request is authenticated
func IsAuthenticated(endpointHandler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		if hasValidToken(response, request) {
			endpointHandler(response, request)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not Authorized"))
		}
	})
}

// VerifyUser handles user with access token
func VerifyUser(endpointHandler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		accessToken := request.Header.Get("Authorization")
		accessToken = strings.Split(accessToken, "Bearer ")[1]

		githubUser, err := GetGithubUser(accessToken)
		if err != nil {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not Authorized"))
			return
		}

		ctx := context.WithValue(request.Context(), "githubUser", *githubUser)
		request = request.WithContext(ctx)
		endpointHandler(response, request)
	})
}
