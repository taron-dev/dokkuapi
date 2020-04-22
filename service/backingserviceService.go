package service

import (
	"github.com/ondro2208/dokkuapi/model"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ServicesService interface {
	CreateService(name string, serviceType string) (*model.Service, int, string)
	GetService(serviceId primitive.ObjectID) (*model.Service, int, string)
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

func (ss *ServicesServiceContext) GetService(serviceId primitive.ObjectID) (*model.Service, int, string) {
	var service = new(model.Service)
	services, ctx := getCollection(ss.store.Client, ss.store.DbName, "services")
	err := services.FindOne(ctx, model.Service{Id: serviceId}).Decode(&service)
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}
	return service, http.StatusOK, "Service founded by id"
}
