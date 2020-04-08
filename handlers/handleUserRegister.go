package handlers

import (
	auth "github.com/ondro2208/dokkuapi/authentication"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// UserRegister handles registration of new user
func UserRegister(w http.ResponseWriter, r *http.Request, store *str.Store) {
	githubUser, err := contextimpl.GetGithubUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	usersService := service.NewUsersService(store)
	user, status, message := usersService.CreateUser(githubUser)
	if user != nil {
		respondAfterVerify(w, r, status, user)
	} else {
		helper.RespondWithMessage(w, r, status, message)
	}
}

func respondAfterVerify(w http.ResponseWriter, r *http.Request, status int, user *model.User) {
	jwt, err := auth.GenerateJWT(user.Id.Hex())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	w.Header().Add("Authorization", "Bearer "+jwt)
	helper.RespondWithData(w, r, status, &user)
}
