package server

import (
	auth "github.com/ondro2208/dokkuapi/authentication"
	author "github.com/ondro2208/dokkuapi/authorization"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

// isAuthenticated verifies if request is authenticated
func (s *Server) isAuthenticated(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		if auth.HasValidToken(response, request, s.blackList) {
			// database usage for user id ??? TODO
			request = contextimpl.DecorateWithSub(request)
			endpointHandler(response, request)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not Authorized"))
		}
	}
}

// verifyUser handles user with access token
func (s *Server) verifyUser(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
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
	}
}

func (s *Server) isUserAuthorizedApp(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sub, err := contextimpl.GetSub(request.Context())
		if err != nil {
			helper.RespondWithMessage(response, request, http.StatusInternalServerError, err.Error())
		}
		responseObj := author.GetAppID(request)
		if responseObj.Value == nil {
			helper.RespondWithMessage(response, request, responseObj.Status, responseObj.Message)
		}
		appID := responseObj.Value.(primitive.ObjectID)

		usersService := service.NewUsersService(s.store)
		user, status, message := usersService.GetExistingUserById(sub)
		if user == nil {
			helper.RespondWithMessage(response, request, status, message)
		}

		responseObj = author.AuthorizeUserApp(user, appID)
		if responseObj.Value == nil {
			helper.RespondWithMessage(response, request, responseObj.Status, responseObj.Message)
		}
		app := responseObj.Value.(*model.Application)

		request = contextimpl.DecorateWithUser(request, user)
		request = contextimpl.DecorateWithApp(request, app)
		endpointHandler(response, request)
	}
}
