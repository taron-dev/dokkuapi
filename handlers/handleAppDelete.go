package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/apps"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// AppDelete delete app after user's authorization
func AppDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	usersService := service.NewUsersService(store)
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := contextimpl.GetUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	//TODO unlink backing services
	// dokku apps:destory
	code, m, err := apps.DestroyApp(app.Name)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, code, m)
		return
	}
	status, message, requireRecreate := usersService.DeleteUserApplication(user.Id, app.Id)
	if requireRecreate {
		log.ErrorLogger.Println(err)
		// destroy created app
		_, _, err := apps.CreateApp(app.Name)
		log.ErrorLogger.Println("Creating already destroyed app was successful ", err == nil)
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	helper.RespondWithMessage(w, r, status, message)
}
