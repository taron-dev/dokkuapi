package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AppDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	//extract app id form url
	appIdParam, present := mux.Vars(r)["appId"]
	if !present || len(appIdParam) < 1 {
		log.ErrorLogger.Println("Extracting appId from url was unsuccessful")
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Unknown app id")
	}
	appId, err := primitive.ObjectIDFromHex(appIdParam)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		helper.RespondWithMessage(w, r, http.StatusNotFound, "Unable parse app id")
		return
	}
	log.GeneralLogger.Println("Extracted application id from url: ", appId)

	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUserById(sub)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
	}

	//Authorize user towards app and get app name
	var appName string = ""
	for _, app := range user.Applications {
		if app.Id.Hex() == appId.Hex() {
			appName = app.Name
			break
		}
	}
	if len(appName) == 0 {
		helper.RespondWithMessage(w, r, http.StatusUnauthorized, "Denied access to application for current user")
		return
	}
	log.GeneralLogger.Println("Founded user in database ", user.Username, " with application ", appName)
	//TODO backing services
	// dokku apps:destory
	err, code, m := plugins.DestroyApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, code, m)
		return
	}
	status, message, requireRecreate := usersService.DeleteUserApplication(user.Id, appId)
	if requireRecreate {
		log.ErrorLogger.Println(err)
		// destroy created app
		err, _, _ := plugins.CreateApp(appName)
		log.ErrorLogger.Println("Creating already destroyed app was successful ", err == nil)
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	helper.RespondWithMessage(w, r, status, message)
}
