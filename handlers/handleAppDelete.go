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

// AppDelete delete app after user's authorization
func AppDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	//extract app id form url
	appIDParam, present := mux.Vars(r)["appId"]
	if !present || len(appIDParam) < 1 {
		log.ErrorLogger.Println("Extracting appId from url was unsuccessful")
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Unknown app id")
	}
	appID, err := primitive.ObjectIDFromHex(appIDParam)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		helper.RespondWithMessage(w, r, http.StatusNotFound, "Unable parse app id")
		return
	}
	log.GeneralLogger.Println("Extracted application id from url: ", appID)

	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUserById(sub)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
	}

	//Authorize user towards app and get app name
	var appName string = ""
	for _, app := range user.Applications {
		if app.Id.Hex() == appID.Hex() {
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
	code, m, err := plugins.DestroyApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, code, m)
		return
	}
	status, message, requireRecreate := usersService.DeleteUserApplication(user.Id, appID)
	if requireRecreate {
		log.ErrorLogger.Println(err)
		// destroy created app
		_, _, err := plugins.CreateApp(appName)
		log.ErrorLogger.Println("Creating already destroyed app was successful ", err == nil)
		helper.RespondWithMessage(w, r, status, message)
		return
	}
	helper.RespondWithMessage(w, r, status, message)
}
