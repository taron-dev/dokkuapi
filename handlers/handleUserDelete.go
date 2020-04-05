package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

func UserDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	//TODO delete related services
	//TODO delete related apps
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	userIdParam := mux.Vars(r)["userId"]
	if sub == userIdParam {
		usersService := service.NewUsersService(store)
		err := usersService.DeleteExistingUser(userIdParam)
		if err != nil {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, "User not deleted")
		} else {
			helper.RespondWithMessage(w, r, http.StatusAccepted, "User deleted")
		}
	} else {
		helper.RespondWithMessage(w, r, http.StatusUnauthorized, "Not Authorized")
	}
}
