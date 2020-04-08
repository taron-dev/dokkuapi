package controller

import (
	auth "github.com/ondro2208/dokkuapi/authentication"
	"github.com/ondro2208/dokkuapi/model"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// PostLogin handles login endpoint
func PostLogin(response http.ResponseWriter, request *http.Request) {
	githubUser := request.Context().Value("githubUser").(model.GithubUser)
	var user model.User

	users, ctx := GetCollection1("users")
	err := users.FindOne(ctx, model.User{GithubId: githubUser.Id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{ "message": "User doesn't exist" }`))
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		}
	} else {
		jwt, err := auth.GenerateJWT(user.Id.Hex())
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		response.Header().Set("Authorization", "Bearer "+jwt)
		response.WriteHeader(http.StatusCreated)
		response.Write([]byte(`{ "username": "` + user.Username + `" ,"message": "Successfully logged in" }`))
	}
}
