package service

import (
	"errors"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type ServicesService interface {
	CreateService(name string, serviceType string) (*model.Service, int, string)
	GetService(serviceId primitive.ObjectID) (*model.Service, int, string)
	GetServiceById(serviceId string) (*model.Service, int, string)
	DeleteExistingService(serviceIdHex string) error
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

func (ss *ServicesServiceContext) GetServiceById(serviceId string) (*model.Service, int, string) {
	idPrimitive, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		log.ErrorLogger.Println("Parsing ObjectId from hex error: ", err.Error())
		return nil, http.StatusInternalServerError, "Can't find user"
	}
	var service = new(model.Service)
	services, ctx := getCollection(ss.store.Client, ss.store.DbName, "services")
	err = services.FindOne(ctx, model.Service{Id: idPrimitive}).Decode(&service)
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}
	return service, http.StatusOK, "Service founded by id"
}

func (us *ServicesServiceContext) DeleteExistingService(serviceIdHex string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(serviceIdHex)
	if err != nil {
		log.ErrorLogger.Println("Parsing ObjectId from hex error: ", err.Error())
		return err
	}
	services, ctx := getCollection(us.store.Client, us.store.DbName, "services")
	res, err := services.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		log.ErrorLogger.Println("Delete one service error: ", err.Error())
		return err
	}
	if res.DeletedCount != 1 {
		message := "Delete " + strconv.FormatInt(res.DeletedCount, 10) + " instead of 1"
		log.ErrorLogger.Println(message)
		return errors.New(message)
	}
	return nil
}
