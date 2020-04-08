package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/plugins"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

func AppsCreate(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	log.GeneralLogger.Println("User id from jwt ", sub)

	var appPost model.ApplicationPost
	err = helper.Decode(w, r, &appPost)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Unable to parse request body")
		return
	}
	appName := appPost.Name
	log.GeneralLogger.Println("Extracted application name ", appName)

	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUserById(sub)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
	}
	log.GeneralLogger.Println("Founded user in database ", user.Username)

	// TODO creating backing services
	// dokku apps:create
	err, code, m := plugins.CreateApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, code, m)
		return
	}
	log.GeneralLogger.Println("Application ", appName, " created successfully")

	// TODO add backing services
	app, status, message := usersService.UpdateUserWithApplication(appName, user.Id)
	if app == nil {
		err, _, _ := plugins.DestroyApp(appName)
		log.ErrorLogger.Println("Destroying already created app was successfull: ", err == nil)
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	helper.RespondWithData(w, r, status, &app)
}
