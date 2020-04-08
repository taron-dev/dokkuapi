package handlers

import (
	auth "github.com/ondro2208/dokkuapi/authentication"
	"github.com/ondro2208/dokkuapi/helper"
	"net/http"
)

// UserLogout blacklist request's jwt
func UserLogout(w http.ResponseWriter, r *http.Request, blackList *[]string) {
	auth.AddToBlacklist(r, blackList)
	helper.RespondWithMessage(w, r, http.StatusCreated, "User is successfully logged out")
}
