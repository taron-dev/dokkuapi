package service

import (
	"context"
	"fmt"
	"github.com/ondro2208/dokkuapi/model"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type UsersService interface {
	CreateUser(githubUser *model.GithubUser) (*model.User, int, string)
}

func NewUsersService(serviceStore *str.Store) UsersService {
	return &UsersServiceContext{store: serviceStore}
}

type UsersServiceContext struct {
	store *str.Store
}

func getCollection(client *mongo.Client, dbName string, collectionName string) (*mongo.Collection, context.Context) {
	collection := client.Database(dbName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return collection, ctx
}

func (us *UsersServiceContext) CreateUser(githubUser *model.GithubUser) (*model.User, int, string) {
	var user = new(model.User)
	users, ctx := getCollection(us.store.Client, us.store.DbName, "users")
	err := users.FindOne(ctx, model.User{GithubId: githubUser.Id}).Decode(&user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, http.StatusInternalServerError, err.Error()
		}
		result, _ := users.InsertOne(ctx, model.User{GithubId: githubUser.Id, Username: githubUser.Login})
		users.FindOne(ctx, model.User{Id: result.InsertedID.(primitive.ObjectID)}).Decode(&user)
		return user, http.StatusCreated, "User created"
	} else {
		return nil, http.StatusConflict, "User already registered"
	}
}
