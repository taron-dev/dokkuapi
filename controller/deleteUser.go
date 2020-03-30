package controller

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/ondro2208/dokkuapi/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strings"
)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	//TODO delete related services
	//TODO delete related apps

	//TODO fix getting sub
	reqToken := request.Header.Get("Authorization")
	reqToken = strings.Split(reqToken, "Bearer ")[1]
	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error parsing jwt")
		}
		return mySigningKey, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)

	userIdParam := mux.Vars(request)["userId"]
	if sub == userIdParam {
		idPrimitive, err := primitive.ObjectIDFromHex(sub)
		if err != nil {
			log.ErrorLogger.Fatal("primitive.ObjectIDFromHex ERROR:", err)
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "User not deleted" }`))
			return
		}
		users, ctx := GetCollection("users")
		res, err := users.DeleteOne(ctx, bson.M{"_id": idPrimitive})
		if err != nil {
			log.ErrorLogger.Fatal("DeleteOne() ERROR:", err)
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "User not deleted" }`))
			return
		}
		if res.DeletedCount != 1 {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "User not deleted" }`))
			return
		}
		response.WriteHeader(http.StatusAccepted)
		response.Write([]byte(`{ "message": "User deleted" }`))
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Not Authorized"))
	}

}
