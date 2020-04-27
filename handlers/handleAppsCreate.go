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

// AppsCreate creates application for user, configured by request body
func AppsCreate(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	log.GeneralLogger.Println("User id from jwt ", sub)

	var appPost postApp
	err = helper.Decode(w, r, &appPost)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Unable to parse request body")
		return
	}
	appName := appPost.Name
	log.GeneralLogger.Println("Extracted application name ", appName)

	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUserById(sub)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	log.GeneralLogger.Println("Founded user in database ", user.Username)

	// dokku apps:create
	code, m, err := apps.CreateApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, code, m)
		return
	}
	log.GeneralLogger.Println("Application ", appName, " created successfully")

	app, status, message := usersService.UpdateUserWithApplication(appName, user.Id)
	if app == nil {
		_, _, err := apps.DestroyApp(appName)
		log.ErrorLogger.Println("Destroying already created app was successfull: ", err == nil)
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	helper.RespondWithData(w, r, status, &app)
}

type postApp struct {
	Name string `json:"appName,omitempty"`
}
