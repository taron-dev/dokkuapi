package controller

import (
	"encoding/json"
	model "github.com/ondro2208/dokkuapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

// RegisterUserEndpoint handles register endpoint
func RegisterUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	accessToken := request.Header.Get("Authorization")
	accessToken = strings.Split(accessToken, "Bearer ")[1]

	var user model.User
	githubUser, err := GetGithubUser(accessToken)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Not Authorized"))
		return
	}

	users, ctx := GetCollection("users")
	err = users.FindOne(ctx, model.User{GithubId: githubUser.Id}).Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		result, _ := users.InsertOne(ctx, model.User{GithubId: githubUser.Id, UserName: githubUser.Login})
		users.FindOne(ctx, model.User{Id: result.InsertedID.(primitive.ObjectID)}).Decode(&user)
		json.NewEncoder(response).Encode(user)
	} else {
		response.WriteHeader(http.StatusConflict)
		response.Write([]byte(`{ "message": "User already registered" }`))
	}

}
