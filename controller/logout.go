package controller

import (
	//auth "github.com/ondro2208/dokkuapi/authentication"
	"net/http"
)

// PostLogout handles logout endpoint
func PostLogout(response http.ResponseWriter, request *http.Request) {
	//auth.AddToBlacklist(request)
	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{ "message": "Successfully logged out" }`))
}
