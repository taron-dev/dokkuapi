package controller

import (
	"encoding/json"
	auth "github.com/ondro2208/dokkuapi/authentication"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/plugins"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// PostApps endpoint create application
func PostApps(response http.ResponseWriter, request *http.Request) {
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

	// extract parameters from body
	var app model.ApplicationPost
	err = json.NewDecoder(request.Body).Decode(&app)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{ "message": "Unable to parse request body" }`))
		return
	}
	appName := app.Name
	log.GeneralLogger.Println("Extracted application name ", appName)

	// TODO creating backing services
	// dokku apps:create
	err, code, m := plugins.CreateApp(appName)
	if err != nil {
		log.ErrorLogger.Println(err)
		response.WriteHeader(code)
		response.Write([]byte(`{ "message": "` + m + `" }`))
		return
	}
	log.GeneralLogger.Println("Application ", appName, " created successfully")

	// load user according jwt
	users, ctx := GetCollection1("users")
	var user model.User
	err = users.FindOne(ctx, model.User{Id: userId}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "Unknown user" }`))
		return
	}
	log.GeneralLogger.Println("Founded user in database ", user.Username)

	// TODO add backing services
	newApp := model.Application{
		Name: appName,
		Id:   primitive.NewObjectID(),
	}
	result, err := users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$push", bson.M{"applications": newApp}},
		},
	)
	log.GeneralLogger.Println("Result after updating database ", result)

	if err != nil {
		log.ErrorLogger.Println(err)
		// destroy created app
		err, _, _ := plugins.DestroyApp(appName)
		log.ErrorLogger.Println("Destroying already created app is", err == nil)
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte(`{ "message": "Unable to store application" }`))
		return
	}

	if result.MatchedCount == 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "Unknown user" }`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{ "message": "` + appName + ` created successfully", "application": { "appName":"` + appName + `","id":"` + newApp.Id.Hex() + `"} }`))

}
