package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// UserLogin handles logging in existing user
func UserLogin(w http.ResponseWriter, r *http.Request, store *str.Store) {
	githubUser, err := contextimpl.GetGithubUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUser(githubUser)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	respondAfterVerify(w, r, status, user)

}
