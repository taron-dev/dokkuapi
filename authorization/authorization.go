package authorization

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// GetAppID provide appId as primitive.ObjectId from url
func GetAppID(request *http.Request) helper.Response {
	appIDParam, present := mux.Vars(request)["appId"]
	if !present || len(appIDParam) < 1 {
		log.ErrorLogger.Println("Extracting appId from url was unsuccessful")
		return helper.Response{Value: nil, Status: http.StatusBadRequest, Message: "Unknown app id"}
	}
	appID, err := primitive.ObjectIDFromHex(appIDParam)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return helper.Response{Value: nil, Status: http.StatusNotFound, Message: "Unknown app id"}
	}
	return helper.Response{Value: appID, Status: http.StatusOK, Message: "App id successfully parsed"}
}

// AuthorizeUserApp authorize app towards user apps
func AuthorizeUserApp(user *model.User, appID primitive.ObjectID) helper.Response {
	var requestApp *model.Application = nil
	for _, app := range user.Applications {
		if app.Id.Hex() == appID.Hex() {
			requestApp = &app
			break
		}
	}
	if requestApp == nil {
		log.ErrorLogger.Println("User unauthorized")
		return helper.Response{Value: nil, Status: http.StatusUnauthorized, Message: "Denied access to application for current user"}
	}
	log.GeneralLogger.Println("User is successfully authorized")
	return helper.Response{Value: requestApp, Status: http.StatusOK, Message: "User is authorized"}
}
