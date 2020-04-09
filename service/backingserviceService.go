package service

import (
	//log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type ServicesService interface {
	CreateService(name string, serviceType string) (*model.Service, int, string)
	// GetExistingUser(githubUser *model.GithubUser) (*model.User, int, string)
	// GetExistingUserById(userIdHex string) (*model.User, int, string)
	// DeleteExistingUser(userIdHex string) error
	// UpdateUserWithApplication(appName string, userId primitive.ObjectID) (*model.Application, int, string)
	// DeleteUserApplication(userId primitive.ObjectID, appId primitive.ObjectID) (int, string, bool)
}

func NewServicesService(serviceStore *str.Store) ServicesService {
	return &ServicesServiceContext{store: serviceStore}
}

type ServicesServiceContext struct {
	store *str.Store
}

func (ss *ServicesServiceContext) CreateService(name string, serviceType string) (*model.Service, int, string) {
	var service = new(model.Service)
	service.Name = name
	service.Type = serviceType
	services, ctx := getCollection(ss.store.Client, ss.store.DbName, "services")
	result, _ := services.InsertOne(ctx, service)
	services.FindOne(ctx, model.Service{Id: result.InsertedID.(primitive.ObjectID)}).Decode(&service)
	return service, http.StatusCreated, "Service created"
}
