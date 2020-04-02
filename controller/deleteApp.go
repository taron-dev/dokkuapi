package controller

import (
	"github.com/gorilla/mux"
	auth "github.com/ondro2208/dokkuapi/authentication"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/plugins"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// DeleteApp endpoint delete specific application
func DeleteApp(response http.ResponseWriter, request *http.Request) {
	// prepare user id object from jwt
	sub := auth.ExtractSub(request)
	log.GeneralLogger.Println("User id from jwt ", sub)
	userId, err := primitive.ObjectIDFromHex(sub)
	if err != nil {
		log.ErrorLogger.Println(err)
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "message": "Unknown user" }`))
		return
	}

	//extract app id form url
	appIdParam, present := mux.Vars(request)["appId"]
	if !present || len(appIdParam) < 1 {
		log.ErrorLogger.Println("Extracting appId from url was unsuccessful")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "message": "Unknown user" }`))
	}
	appId, err := primitive.ObjectIDFromHex(appIdParam)
	if err != nil {
		log.ErrorLogger.Fatal(err)
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "message": "Unable parse app id" }`))
		return
	}

	log.GeneralLogger.Println("Extracted application id from url: ", appId)

	// get user according jwt
	users, ctx := GetCollection("users")
	var user model.User
	//err = users.FindOne(ctx, model.User{Id: userId, Applications}).Decode(&user)
	err = users.FindOne(ctx, model.User{Id: userId}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "message": "Unknown user" }`))
		return
	}
	log.GeneralLogger.Println("Founded user in database ", user.Username, " with application id ", appId)

	//Authorize user towards app and get app name
	var appName string = ""
	for _, app := range user.Applications {
		if app.Id.Hex() == appId.Hex() {
			appName = app.Name
			break
		}
	}
	if len(appName) == 0 {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{ "message": "Denied access to application for current user" }`))
	}
	log.GeneralLogger.Println("Founded application name: ", appName)

	//TODO backing services

	// dokku apps:destory
	err, code, m := plugins.DestroyApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		response.WriteHeader(code)
		response.Write([]byte(`{ "message": "` + m + `" }`))
		return
	}

	// delete destroyed app from user's applications
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$pull", bson.M{"applications": bson.M{"_id": appId}}},
		},
	)
	log.GeneralLogger.Println("Result after deleting application ", result)

	if err != nil {
		log.ErrorLogger.Println(err)
		// destroy created app
		err, _, _ := plugins.CreateApp(appName)
		log.ErrorLogger.Println("Creating already destroyed app is successful ", err == nil)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "Unable to delete application" }`))
		return
	}

	if result.MatchedCount == 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "No application deleted" }`))
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{ "message": "` + appName + ` destroyed successfully" }`))

}
