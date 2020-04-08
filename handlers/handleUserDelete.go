package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// UserDelete delete authorized user from db
func UserDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	//TODO delete related services
	//TODO delete related apps
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	userIDParam := mux.Vars(r)["userId"]
	if sub == userIDParam {
		usersService := service.NewUsersService(store)
		err := usersService.DeleteExistingUser(userIDParam)
		if err != nil {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, "User not deleted")
		} else {
			helper.RespondWithMessage(w, r, http.StatusAccepted, "User deleted")
		}
	} else {
		helper.RespondWithMessage(w, r, http.StatusUnauthorized, "Not Authorized")
	}
}
