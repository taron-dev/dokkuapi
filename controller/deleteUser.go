package controller

import (
	"github.com/gorilla/mux"
	auth "github.com/ondro2208/dokkuapi/authentication"
	log "github.com/ondro2208/dokkuapi/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	//TODO delete related services
	//TODO delete related apps

	//TODO fix getting sub
	sub := auth.ExtractSub(request)
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
